package main

import(
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	"path/filepath"
	"k8s.io/client-go/util/homedir"
	"log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"github.com/appscode/envconfig"
	"github.com/appscode/go-notify/unified"
	"strings"
	"github.com/appscode/go-notify"
	"fmt"
)

type Receiver struct {
	// To whom notification will be sent
	To []string `json:"to,omitempty"`

	// How this notification will be sent
	Notifier string `json:"notifier,omitempty"`
}

func main() {
	receiver := Receiver{Notifier:"Hipchat", To:[]string{"ops-alerts"}}
	//receiver := Receiver{Notifier:"Twilio", To:[]string{"+8801845683020"}}
	kubeConfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	secret, err := kubeClient.CoreV1().
		Secrets("kube-system").
		Get("notifier-config", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	notifierCred := func(key string) (value string, found bool) {
		var bytes []byte
		bytes, found = secret.Data[key]
		value = string(bytes)
		return
	}

	fmt.Println("loading secret data......")

	notifier, err := unified.LoadVia(strings.ToLower(receiver.Notifier), notifierCred)
	if err != nil {
		log.Fatal(err)
	}

	switch n := notifier.(type) {
	case notify.ByEmail:
		err = fmt.Errorf("expected hipchat notification")
	case notify.BySMS:
		err = n.To(receiver.To[0], receiver.To[1:]...).
			WithBody("Hello world. Checking whether notifications are being send or not").
			Send()
	case notify.ByChat:
		err = n.To(receiver.To[0], receiver.To[1:]...).
			WithBody("Hello world. Checking whether notifications are being send or not").
			Send()
	case notify.ByPush:
		err = fmt.Errorf("expected hipchat notification")
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("success.....")
}