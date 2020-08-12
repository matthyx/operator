package cautils

import (
	"fmt"
	"strings"
)

var (
	WlidPrefix          = "wlid://"
	ClusterWlidPrefix   = "cluster-"
	NamespaceWlidPrefix = "namespace-"
	K8SKindsList        = []string{"ComponentStatus", "ConfigMap", "ControllerRevision", "CronJob",
		"CustomResourceDefinition", "DaemonSet", "Deployment", "Endpoints", "Event", "HorizontalPodAutoscaler",
		"Ingress", "Job", "Lease", "LimitRange", "LocalSubjectAccessReview", "MutatingWebhookConfiguration",
		"Namespace", "NetworkPolicy", "Node", "PersistentVolume", "PersistentVolumeClaim", "Pod",
		"PodDisruptionBudget", "PodSecurityPolicy", "PodTemplate", "PriorityClass", "ReplicaSet",
		"ReplicationController", "ResourceQuota", "Role", "RoleBinding", "Secret", "SelfSubjectAccessReview",
		"SelfSubjectRulesReview", "Service", "ServiceAccount", "StatefulSet", "StorageClass",
		"SubjectAccessReview", "TokenReview", "ValidatingWebhookConfiguration", "VolumeAttachment"}
	KindReverseMap = map[string]string{}
)

//
func restoreInnerIdentifiersFromWLID(spiffeSlices []string) []string {
	for _, kind := range K8SKindsList {
		KindReverseMap[strings.ToLower(strings.Replace(kind, "-", "", -1))] = kind
	}
	if strings.HasPrefix(spiffeSlices[0], ClusterWlidPrefix) &&
		strings.HasPrefix(spiffeSlices[1], NamespaceWlidPrefix) &&
		strings.Contains(spiffeSlices[2], "-") {
		spiffeSlices[0] = spiffeSlices[0][len(ClusterWlidPrefix):]
		spiffeSlices[1] = spiffeSlices[1][len(NamespaceWlidPrefix):]
		dashIdx := strings.Index(spiffeSlices[2], "-")
		spiffeSlices = append(spiffeSlices, spiffeSlices[2][dashIdx+1:])
		spiffeSlices[2] = spiffeSlices[2][:dashIdx]
		if val, ok := KindReverseMap[spiffeSlices[2]]; ok {
			spiffeSlices[2] = val
		}
	}

	return spiffeSlices
}

// RestoreMicroserviceIDsFromSpiffe -
func RestoreMicroserviceIDsFromSpiffe(spiffe string) ([]string, error) {
	if StringHasWhitespace(spiffe) {
		return nil, fmt.Errorf("wlid %s invalid. whitespace found", spiffe)
	}
	if strings.HasPrefix(spiffe, WlidPrefix) {
		spiffe = spiffe[len(WlidPrefix):]
	}
	spiffeSlices := strings.Split(spiffe, "/")
	// The documented WLID format (https://cyberarmorio.sharepoint.com/sites/development2/Shared%20Documents/kubernetes_design1.docx?web=1)
	if len(spiffeSlices) == 3 {
		spiffeSlices = restoreInnerIdentifiersFromWLID(spiffeSlices)
	}
	if len(spiffeSlices) != 4 { // first used WLID, deprecated since 24.10.2019
		return spiffeSlices, fmt.Errorf("invalid WLID format")
	}
	return spiffeSlices, nil
}

// StringHasWhitespace check if a string has whitespace
func StringHasWhitespace(str string) bool {
	if whitespace := strings.Index(str, " "); whitespace != -1 {
		return true
	}

	return false
}

// GetK8SKindFronList get the calculated wlid
func GetK8SKindFronList(kind string) string {
	for i := range K8SKindsList {
		if strings.ToLower(kind) == strings.ToLower(K8SKindsList[i]) {
			return K8SKindsList[i]
		}
	}
	return kind
}

// GetNamespaceFromWlid parse wlid and get Namespace
func GetNamespaceFromWlid(wlid string) string {
	r, err := RestoreMicroserviceIDsFromSpiffe(wlid)
	if err != nil {
		return ""
	}
	return r[1]
}

// GetKindFromWlid parse wlid and get kind
func GetKindFromWlid(wlid string) string {
	r, err := RestoreMicroserviceIDsFromSpiffe(wlid)
	if err != nil {
		return ""
	}
	return GetK8SKindFronList(r[2])
}

// GetNameFromWlid parse wlid and get name
func GetNameFromWlid(wlid string) string {
	r, err := RestoreMicroserviceIDsFromSpiffe(wlid)
	if err != nil {
		return ""
	}
	return r[3]
}

// GetWLID get the calculated wlid
func GetWLID(level0, level1, k, name string) string {
	kind := strings.ToLower(k)
	kind = strings.Replace(kind, "-", "", -1)
	return fmt.Sprintf("%s%s%s/%s%s/%s-%s", WlidPrefix, ClusterWlidPrefix, level0, NamespaceWlidPrefix, level1, kind, name)

}

// IsWalidValid test if wlid is a valid wlid
func IsWalidValid(wlid string) error {
	_, err := RestoreMicroserviceIDsFromSpiffe(wlid)
	return err
}