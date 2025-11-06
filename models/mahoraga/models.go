package mahoraga

import "time"

type App struct {
	Id          string     `json:"id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description,omitempty"`
	TenantId    *string    `json:"tenant_id,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type AppSettings struct {
	Id            *int       `json:"id"`
	AppId         *string    `json:"app_id"`
	MaxStorageMB  *float32   `json:"max_storage_mb,omitempty"`
	UsedStorageMB *float32   `json:"used_storage_mb,omitempty"`
	MaxFileSizeMB *float32   `json:"max_file_size_mb,omitempty"`
	State         *string    `json:"state,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

type File struct {
	Id         int64     `json:"id"`
	FileId     string    `json:"file_id"`
	AppId      string    `json:"app_id"`
	Filename   string    `json:"file_name"`
	Type       string    `json:"type"`
	FileSizeMB float32   `json:"file_size"`
	Deleted    bool      `json:"deleted"`
	Content    []byte    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FilesInfo struct {
	FileId      string    `json:"file_id"`
	DownloadURL string    `json:"download_url"`
	Filename    string    `json:"file_name"`
	Type        string    `json:"file_type"`
	FileSizeMB  float32   `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
}

type User struct {
	Id        string  `json:"id"`
	Username  *string `json:"username,omitempty"`
	Password  *string `json:"password,omitempty"`
	Role      *string `json:"role,omitempty"`
	Active    *bool   `json:"active,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdatedAt *string `json:"updated_at,omitempty"`
}

type Login struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=3"`
}

type LoginServicePayload struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Payload[T any] struct {
	Content T        `json:"content"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

type MahoragaResponse[T any] struct {
	Code     int        `json:"code"`
	HasError bool       `json:"hasError"`
	Payload  Payload[T] `json:"payload"`
}

type MahhoragaPaginatedResponse[T any] struct {
	Code         int        `json:"code"`
	HasError     bool       `json:"hasError"`
	Payload      Payload[T] `json:"payload"`
	Page         int        `json:"page"`
	Limit        int        `json:"limit"`
	TotalResults int        `json:"totalResults"`
	Pages        int        `json:"pages"`
}
