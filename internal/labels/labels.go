package labels

import (
	"fmt"
	"strings"
)

const (
	Prefix = "io.balena.image."

	// Class identifies the extension type. Required.
	Class = Prefix + "class"

	// ClassOverlay is the only supported class value.
	ClassOverlay = "overlay"

	// RequiresReboot indicates the host must reboot after install. Optional.
	RequiresReboot = Prefix + "requires-reboot"

	// KernelVersion declares kernel ABI version (M.m.p) for userspace compatibility. Optional.
	KernelVersion = Prefix + "kernel-version"

	// KernelABIID identifies the kernel's binary interface for module/eBPF compatibility. Optional.
	KernelABIID = Prefix + "kernel-abi-id"
)

// Validate checks that the OCI annotations contain the required extension labels.
func Validate(annotations map[string]string) error {
	class, ok := annotations[Class]
	if !ok {
		return fmt.Errorf("missing required label %s", Class)
	}
	if class != ClassOverlay {
		return fmt.Errorf("unsupported %s=%q, must be %q", Class, class, ClassOverlay)
	}
	return nil
}

// ToEnv converts io.balena.image.* annotations to environment variables.
// "io.balena.image.class" becomes "EXTENSION_IMAGE_CLASS=overlay".
func ToEnv(annotations map[string]string) []string {
	var env []string
	for k, v := range annotations {
		if !strings.HasPrefix(k, Prefix) {
			continue
		}
		suffix := strings.TrimPrefix(k, Prefix)
		name := "EXTENSION_IMAGE_" + strings.ToUpper(strings.ReplaceAll(suffix, "-", "_"))
		env = append(env, name+"="+v)
	}
	return env
}
