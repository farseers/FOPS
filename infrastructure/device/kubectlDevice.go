package device

import (
	"context"
	"fops/domain/_/eumK8SControllers"
	"fops/domain/apps"
	"fops/domain/cluster"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
)

func RegisterKubectlDevice() {
	container.Register(func() apps.IKubectlDevice { return &kubectlDevice{} })
}

type kubectlDevice struct {
}

func (kubectlDevice) GetConfigFile(clusterName string) string {
	return apps.KubeRoot + clusterName
}

func (device kubectlDevice) CreateConfigFile(clusterName string, clusterConfig string) string {
	configFile := device.GetConfigFile(clusterName)
	// 文件不存在，则创建
	if !file.IsExists(configFile) {
		file.WriteString(configFile, clusterConfig)
	} else {
		// 比对配置是否不一样，不一样则覆盖新的
		var config = file.ReadString(configFile)
		if clusterConfig != config {
			file.WriteString(configFile, clusterConfig)
		}
	}
	return configFile
}

func (device kubectlDevice) SetYaml(clusterName string, projectName string, yamlContent string, progress chan string, ctx context.Context) bool {
	// 将yaml文件写入临时文件
	fileName := "/tmp/" + projectName + ".yaml"
	file.Delete(fileName)
	file.WriteString(fileName, yamlContent)

	configFile := device.GetConfigFile(clusterName)

	// 发布
	lstResult, wait := exec.RunShellContext(ctx, "kubectl apply -f "+fileName+" --kubeconfig="+configFile+" --insecure-skip-tls-verify", nil, "", true)
	if exitCode := exec.SaveToChan(progress, lstResult, wait); exitCode != 0 {
		progress <- "K8S更新镜像失败。"
		return false
	}
	progress <- "更新镜像版本完成。"
	return true
}

func (device kubectlDevice) SetImages(cluster cluster.DomainObject, projectName string, dockerImages string, k8SControllersType eumK8SControllers.Enum, progress chan string, ctx context.Context) bool {
	return device.SetImagesByClusterName(cluster.Name, "cluster.K8sConnectConfig", projectName, dockerImages, k8SControllersType, progress, ctx)

}

func (device kubectlDevice) SetImagesByClusterName(clusterName string, clusterConfig string, projectName string, dockerImages string, k8SControllersType eumK8SControllers.Enum, progress chan string, ctx context.Context) bool {

	progress <- "---------------------------------------------------------"
	progress <- "开始更新K8S POD的镜像版本。"

	var configFile = device.CreateConfigFile(clusterName, clusterConfig)
	lstResult, wait := exec.RunShellContext(ctx, "kubectl set image "+k8SControllersType.String()+"/"+projectName+" "+projectName+"="+dockerImages+" --kubeconfig="+configFile+" --insecure-skip-tls-verify", nil, "", true)
	if exitCode := exec.SaveToChan(progress, lstResult, wait); exitCode != 0 {
		progress <- "K8S更新镜像失败。"
		return false
	}
	progress <- "更新镜像版本完成。"
	return true
}
