package export

import (
	"github.com/bitrise-tools/go-xcode/certificateutil"
	"github.com/bitrise-tools/go-xcode/profileutil"
)

// CodeSignGroup ...
type CodeSignGroup interface {
	GetCertificate() certificateutil.CertificateInfoModel
	GetInstallerCertificate() *certificateutil.CertificateInfoModel
	GetBundleIDProfileMap() map[string]profileutil.ProvisioningProfileInfoModel
}
