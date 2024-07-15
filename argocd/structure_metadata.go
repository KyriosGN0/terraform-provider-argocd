package argocd

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func expandMetadata(d *schema.ResourceData) (meta meta.ObjectMeta) {
	m := d.Get("metadata.0").(map[string]interface{})

	if v, ok := m["annotations"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Annotations = expandStringMap(m["annotations"].(map[string]interface{}))
	}

	if v, ok := m["labels"].(map[string]interface{}); ok && len(v) > 0 {
		meta.Labels = expandStringMap(m["labels"].(map[string]interface{}))
	}

	if v, ok := m["name"]; ok {
		meta.Name = v.(string)
	}

	if v, ok := m["namespace"]; ok {
		meta.Namespace = v.(string)
	}

	if v, ok := m["finalizers"].([]interface{}); ok && len(v) > 0 {
		meta.Finalizers = expandStringList(v)
	}

	return meta
}

func flattenMetadata(meta meta.ObjectMeta, d *schema.ResourceData) []interface{} {
	m := map[string]interface{}{
		"generation":       meta.Generation,
		"name":             meta.Name,
		"namespace":        meta.Namespace,
		"resource_version": meta.ResourceVersion,
		"uid":              fmt.Sprintf("%v", meta.UID),
	}

	annotations := d.Get("metadata.0.annotations").(map[string]interface{})
	m["annotations"] = metadataRemoveInternalKeys(meta.Annotations, annotations)

	labels := d.Get("metadata.0.labels").(map[string]interface{})
	m["labels"] = metadataRemoveInternalKeys(meta.Labels, labels)

	return []interface{}{m}
}

func metadataRemoveInternalKeys(m map[string]string, d map[string]interface{}) map[string]string {
	for k := range m {
		if metadataIsInternalKey(k) && !isKeyInMap(k, d) {
			delete(m, k)
		}
	}

	return m
}

func metadataIsInternalKey(annotationKey string) bool {
	u, err := url.Parse("//" + annotationKey)
	if err != nil {
		return false
	}

	return strings.HasSuffix(u.Hostname(), "kubernetes.io") || annotationKey == "notified.notifications.argoproj.io"
}
