package kuberclient

import (
	"k8s_rbac/pkg/logger"
	"os"

	"context"
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/mitchellh/go-homedir"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeConfig struct {
	APIVersion     string    `yaml:"apiVersion"`
	Kind           string    `yaml:"kind"`
	Clusters       []Cluster `yaml:"clusters"`
	Contexts       []Context `yaml:"contexts"`
	CurrentContext string    `yaml:"current-context"`
	Users          []User    `yaml:"users"`
}

type Cluster struct {
	Name    string      `yaml:"name"`
	Cluster ClusterData `yaml:"cluster"`
}

type ClusterData struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type Context struct {
	Name    string      `yaml:"name"`
	Context ContextData `yaml:"context"`
}

type ContextData struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type User struct {
	Name string   `yaml:"name"`
	User UserData `yaml:"user"`
}

type UserData struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}

func FindKubeConfig() (string, error) {
	env := os.Getenv("KUBECONFIG")
	if env != "" {
		return env, nil
	}
	path, err := homedir.Expand("~/.kube/config")
	if err != nil {
		return "", err
	}
	return path, nil
}

func createClientSet(kubeconfigPath string) *kubernetes.Clientset {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		logger.ErrorLogger.Println("Error building kubeconfig", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.ErrorLogger.Println("Error creating clientset", err)
	}

	return clientset
}

func CreateRoles(username string, namespace string) {

	userRoleName := fmt.Sprintf("%s-role", username)
	userRoleBindingName := fmt.Sprintf("%s-rolebinding", username)

	kubeconfigPath, err := FindKubeConfig()
	if err != nil {
		logger.ErrorLogger.Println("Error finding kubeconfig path", err)
		return
	}
	clientset := createClientSet(kubeconfigPath)

	userRole := &rbacv1.Role{
		ObjectMeta: v1.ObjectMeta{
			Name:      userRoleName,
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "services", "endpoints", "events", "pods/log", "pods/portforward", "secrets"},
				Verbs:     []string{"get", "list", "watch", "create", "delete"},
			},
			{
				APIGroups: []string{"apps"},
				Resources: []string{"deployments"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}

	userRoleBinding := &rbacv1.RoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("%s-rolebinding", username),
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind: "User",
				Name: username,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     fmt.Sprintf("%s-role", username),
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	if _, err := clientset.RbacV1().Roles(namespace).Get(context.Background(), userRoleName, v1.GetOptions{}); errors.IsNotFound(err) {
		if _, err := clientset.RbacV1().Roles(namespace).Create(context.Background(), userRole, v1.CreateOptions{}); err != nil {
			logger.ErrorLogger.Println("Error creating Role", err)
		} else {
			logger.InfoLogger.Printf("Role %s created successfully!", userRoleName)
		}
	} else {
		logger.InfoLogger.Printf("Role %s already exists", userRoleName)
	}

	if _, err := clientset.RbacV1().RoleBindings(namespace).Get(context.Background(), userRoleBindingName, v1.GetOptions{}); errors.IsNotFound(err) {
		if _, err := clientset.RbacV1().RoleBindings(namespace).Create(context.Background(), userRoleBinding, v1.CreateOptions{}); err != nil {
			logger.ErrorLogger.Println("Error creating RoleBinding", err)
		} else {
			logger.InfoLogger.Printf("RoleBinding %s created successfully!", userRoleBindingName)
		}
	} else {
		logger.InfoLogger.Printf("RoleBinding %s already exists", userRoleBindingName)
	}
}

func CreateKubeConfig(username string, caB64 string, certB64 string, keyB64 string) {

	kubeconfig := KubeConfig{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentContext: "kind-kind",
		Clusters: []Cluster{
			{
				Name: "kind-cluster",
				Cluster: ClusterData{
					CertificateAuthorityData: caB64,
					Server:                   "https://127.0.0.1:41581",
				},
			},
		},
		Contexts: []Context{
			{
				Name: "kind-kind",
				Context: ContextData{
					Cluster: "kind-cluster",
					User:    username,
				},
			},
		},
		Users: []User{
			{
				Name: username,
				User: UserData{
					ClientCertificateData: certB64,
					ClientKeyData:         keyB64,
				},
			},
		},
	}

	yamlBytes, err := yaml.Marshal(kubeconfig)
	if err != nil {
		logger.ErrorLogger.Println("Error: Failed to marshal yaml ", err)
	}

	err = os.WriteFile("config", yamlBytes, 0644)
	if err != nil {
		logger.ErrorLogger.Println("Error: Failed to write kubeconfig file", err)
	}
}
