package services

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "errors"
    "time"
    "notes-app/config"
    "notes-app/models"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

func NewAuthService() *AuthService {
    return &AuthService{}
}

func (s *AuthService) Register(req models.RegisterRequest) error {
    ctx := context.Background()
    
    // Check if email exists
    var existingUser models.User
    err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
    if err == nil {
        return errors.New("email already exists")
    }

    // Check if username exists
    err = config.DB.Collection("users").FindOne(ctx, bson.M{"username": req.Username}).Decode(&existingUser)
    if err == nil {
        return errors.New("username already exists")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := models.User{
        ID:       primitive.NewObjectID(),
        Name:     req.Name,
        Email:    req.Email,
        Username: req.Username,
        Password: string(hashedPassword),
    }

    _, err = config.DB.Collection("users").InsertOne(ctx, user)
    return err
}

func (s *AuthService) Login(req models.LoginRequest) (models.LoginResponse, error) {
    ctx := context.Background()
    
    var user models.User
    err := config.DB.Collection("users").FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
    if err != nil {
        return models.LoginResponse{}, errors.New("invalid credentials")
    }

    // Verify password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return models.LoginResponse{}, errors.New("invalid credentials")
    }

    // Generate session ID
    sessionID := generateSessionID()
    
    // Store session in Redis (expires in 24 hours)
    err = config.RedisClient.Set(ctx, "session:"+sessionID, user.ID.Hex(), 24*time.Hour).Err()
    if err != nil {
        return models.LoginResponse{}, err
    }

    response := models.LoginResponse{
        SessionID: sessionID,
        User: models.UserProfileDto{
            ID:       user.ID.Hex(),
            Username: user.Username,
            Email:    user.Email,
        },
    }

    return response, nil
}

func (s *AuthService) Logout(sessionID string) error {
    ctx := context.Background()
    return config.RedisClient.Del(ctx, "session:"+sessionID).Err()
}

func (s *AuthService) GetUserBySession(sessionID string) (*models.User, error) {
    ctx := context.Background()
    
    userID, err := config.RedisClient.Get(ctx, "session:"+sessionID).Result()
    if err != nil {
        return nil, errors.New("invalid session")
    }

    objID, err := primitive.ObjectIDFromHex(userID)
    if err != nil {
        return nil, err
    }

    var user models.User
    err = config.DB.Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (s *AuthService) ChangePassword(userID primitive.ObjectID, req models.ChangePasswordRequest) error {
    ctx := context.Background()
    
    var user models.User
    err := config.DB.Collection("users").FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
    if err != nil {
        return err
    }

    // Verify old password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
    if err != nil {
        return errors.New("old password is incorrect")
    }

    // Hash new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // Update password
    _, err = config.DB.Collection("users").UpdateOne(
        ctx,
        bson.M{"_id": userID},
        bson.M{"$set": bson.M{"password": string(hashedPassword)}},
    )

    return err
}

func generateSessionID() string {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}
