package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Post represents a dog's social media post
type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DogID     uuid.UUID `gorm:"type:uuid;not null;index" json:"dog_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	ImageURL  string    `gorm:"type:varchar(255)" json:"image_url"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"created_at"`

	// Relationships
	Dog       Dog         `gorm:"foreignKey:DogID;constraint:OnDelete:CASCADE" json:"dog,omitempty"`
	Likes     []Like      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"likes,omitempty"`
	Comments  []Comment   `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
	Hashtags  []Hashtag   `gorm:"many2many:post_hashtags;constraint:OnDelete:CASCADE" json:"hashtags,omitempty"`
}

// BeforeCreate sets the ID before creating the post
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the Post model
func (Post) TableName() string {
	return "posts"
}

// Like represents a like on a post
type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;index" json:"post_id"`
	DogID     uuid.UUID `gorm:"type:uuid;not null;index" json:"dog_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"post,omitempty"`
	Dog  Dog  `gorm:"foreignKey:DogID;constraint:OnDelete:CASCADE" json:"dog,omitempty"`
}

// BeforeCreate sets the ID before creating the like
func (l *Like) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the Like model
func (Like) TableName() string {
	return "likes"
}

// Comment represents a comment on a post
type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null;index" json:"post_id"`
	DogID     uuid.UUID `gorm:"type:uuid;not null;index" json:"dog_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Post Post `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"post,omitempty"`
	Dog  Dog  `gorm:"foreignKey:DogID;constraint:OnDelete:CASCADE" json:"dog,omitempty"`
}

// BeforeCreate sets the ID before creating the comment
func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the Comment model
func (Comment) TableName() string {
	return "comments"
}

// Hashtag represents a hashtag that can be used in posts
type Hashtag struct {
	ID   uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Tag  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"tag"`

	// Relationships
	Posts []Post `gorm:"many2many:post_hashtags;constraint:OnDelete:CASCADE" json:"posts,omitempty"`
}

// BeforeCreate sets the ID before creating the hashtag
func (h *Hashtag) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

// TableName returns the table name for the Hashtag model
func (Hashtag) TableName() string {
	return "hashtags"
}

// PostHashtag represents the many-to-many relationship between posts and hashtags
type PostHashtag struct {
	PostID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"post_id"`
	HashtagID uuid.UUID `gorm:"type:uuid;primaryKey" json:"hashtag_id"`
}

// TableName returns the table name for the PostHashtag model
func (PostHashtag) TableName() string {
	return "post_hashtags"
}

// Follower represents the follower relationship between dogs
type Follower struct {
	FollowerDogID uuid.UUID `gorm:"type:uuid;primaryKey" json:"follower_dog_id"`
	FollowedDogID uuid.UUID `gorm:"type:uuid;primaryKey" json:"followed_dog_id"`

	// Relationships
	FollowerDog Dog `gorm:"foreignKey:FollowerDogID;constraint:OnDelete:CASCADE" json:"follower_dog,omitempty"`
	FollowedDog Dog `gorm:"foreignKey:FollowedDogID;constraint:OnDelete:CASCADE" json:"followed_dog,omitempty"`
}

// TableName returns the table name for the Follower model
func (Follower) TableName() string {
	return "followers"
}