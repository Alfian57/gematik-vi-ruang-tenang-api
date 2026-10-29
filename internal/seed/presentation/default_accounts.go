package presentation

import (
	"errors"
	"time"

	"github.com/Alfian57/ruang-tenang-api/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	presentationAdminEmail          = "admin@gmail.com"
	presentationAdminPassword       = "password"
	presentationAdminWhatsApp       = "+62895363116178"
	presentationMitraEmail          = "mitra@gmail.com"
	presentationGadingEmail         = "gading@gmail.com"
	presentationDeryEmail           = "dery@gmail.com"
	presentationAndhikaEmail        = "andhika@gmail.com"
	presentationDefaultProfileTheme = "default"

	// Moderation-state showcase accounts (varied banned/blocked/suspended states
	// for testing the admin & moderation pages).
	presentationBlockedEmail     = "blocked.user@ruang-tenang.com"
	presentationBannedEmail      = "banned.user@ruang-tenang.com"
	presentationSuspendedEmail   = "suspended.user@ruang-tenang.com"
	presentationActiveExtraEmail = "active.extra@ruang-tenang.com"
)

var presentationAccountEmails = []string{
	presentationAdminEmail,
	presentationMitraEmail,
	presentationGadingEmail,
	presentationDeryEmail,
	presentationAndhikaEmail,
	presentationBlockedEmail,
	presentationBannedEmail,
	presentationSuspendedEmail,
	presentationActiveExtraEmail,
}

// SeedDefaultAccounts seeds the single fixed admin account for presentation.
func SeedDefaultAccounts(db *gorm.DB) error {
	if err := cleanupPresentationAccounts(db); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(presentationAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	admin := model.User{
		Name:               "Admin",
		Email:              presentationAdminEmail,
		WhatsAppNumber:     presentationAdminWhatsApp,
		WhatsAppVerifiedAt: &now,
		Password:           string(hashedPassword),
		Role:               model.RoleAdmin,
		Exp:                0,
		ProfileTheme:       presentationDefaultProfileTheme,
	}

	var existing model.User
	err = db.Unscoped().Where("email = ?", admin.Email).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Create(&admin).Error
		}
		return err
	}

	return db.Unscoped().Model(&existing).Updates(map[string]interface{}{
		"name":                 admin.Name,
		"whatsapp_number":      presentationAdminWhatsApp,
		"whatsapp_verified_at": now,
		"password":             admin.Password,
		"role":                 admin.Role,
		"exp":                 admin.Exp,
		"is_premium":          false,
		"premium_since":       nil,
		"premium_expires_at":  nil,
		"gold_coins":          0,
		"is_blocked":          false,
		"is_forum_blocked":    false,
		"is_banned":           false,
		"suspension_end":      nil,
		"suspension_reason":   "",
		"ban_reason":          "",
		"profile_theme":       presentationDefaultProfileTheme,
		"profile_banner":      "",
		"avatar_border_color": "",
		"tagline":             "",
		"bio":                 "",
		"deleted_at":          nil,
	}).Error
}

func cleanupPresentationAccounts(db *gorm.DB) error {
	obsoleteEmails := []string{
		"admin@ruang-tenang.com",
		"mitra@ruang-tenang.com",
		"demo.utama@ruang-tenang.com",
		"demo.cadangan@ruang-tenang.com",
		"usertest@gmail.com",
		"mitra.invited@ruang-tenang.com",
	}
	if err := db.Where("email IN ?", obsoleteEmails).Delete(&model.User{}).Error; err != nil {
		return err
	}

	// Never prune accounts based on role. The presentation seeder can be pointed
	// at a database that also contains real users, so cleanup must stay limited
	// to the exact legacy demo identities listed above.
	return nil
}
