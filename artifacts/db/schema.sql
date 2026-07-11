
CREATE TABLE user (
    -- ID System Assigned
    id INT AUTO_INCREMENT PRIMARY KEY,
    
    -- UserFullName (User Assigned - Idx / NOT UNIQUE)
    user_full_name VARCHAR(100) NOT NULL,
    
    -- Email address (User Assigned - Idx & Unique)
    email VARCHAR(100) NOT NULL UNIQUE,
    
    -- Password (System Assigned - VARCHAR(100))
    password VARCHAR(100) NOT NULL,
    
    -- Phone number (User Assigned)
    phone VARCHAR(20) NULL,
    
    -- Update password (System Assigned)
    requires_password_update BOOLEAN NOT NULL DEFAULT FALSE,
    
    -- Role (User Assigned - Idx)
    role ENUM('Super_Admin', 'Admin', 'User') NOT NULL DEFAULT 'User',
    
    -- Creation date (System Assigned)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Update date (System Assigned)
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Status (User Assigned - Idx)
    status ENUM('Active', 'Inactive') NOT NULL DEFAULT 'Active',
    
    -- Change history (System Assigned)
    change_history JSON NULL,

    -- INDEXES
    INDEX idx_user_username (username),
    INDEX idx_user_role (role),
    INDEX idx_user_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
