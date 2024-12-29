package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	metav1 "iam/pkg/api/meta/v1"
	"iam/pkg/util/idutil"
	"iam/pkg/validation"
	"iam/pkg/validation/field"
	"time"
)

type User struct {
	metav1.ObjectMeta `json:"metadata,omitempty"` // 通用
	NickName          string                      `json:"nickname" gorm:"column:nickname"`
	Status            int                         `json:"status" grom:"column:status"`                         // user 状态, 1 表示正常
	Password          string                      `json:"password" gorm:"column:password" validate:"required"` // 标签来确保字段的值不为空
	LoginedAt         *time.Time                  `json:"loginedAt,omitempty" gorm:"column:loginedAt"`
	IsAdmin           int                         `json:"isAdmin,omitempty" gorm:"column:isAdmin" validate:"omitempty"` // 某些字段为空时不在JSON中显示。这时可以使用omitempty标签来实现这一功能。omitempty标签的作用是当字段的值为空时，不将该字段包含在JSON中。
	Email             string                      `json:"email" gorm:"column:email" validate:"required,email,min=1,max=100"`
	// DeletedAt         time.Time                   `json:"deletedAt"`
}

func (u *User) TableName() string {
	return "user"
}

type UserList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*User `json:"items"`
}

// Compare 验证密码是否匹配
func (u *User) Compare(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
}

// AfterCreate 创建新数据后，进行添加 InstanceID
func (u *User) AfterCreate(tx *gorm.DB) error {
	u.InstanceID = idutil.GetInstanceID(u.ID, "user-")

	return tx.Save(u).Error
}

// Validate 验证用户对象是否有效
func (u *User) Validate() field.ErrorList {
	val := validation.NewValidator(u)
	allErrs := val.Validate()

	if err := validation.IsValidPassword(u.Password); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("password"), err.Error(), ""))
	}
	return allErrs
}
