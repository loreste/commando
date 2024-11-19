-- Create the database
CREATE DATABASE IF NOT EXISTS teltech_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use the database
USE teltech_db;

-- Create the users table
CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(100) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       role ENUM('user', 'admin') NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create the folders table
CREATE TABLE folders (
                         id INT AUTO_INCREMENT PRIMARY KEY,
                         name VARCHAR(255) NOT NULL,
                         path VARCHAR(255) UNIQUE NOT NULL,
                         owner_id INT NOT NULL,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create the permissions table
CREATE TABLE permissions (
                             id INT AUTO_INCREMENT PRIMARY KEY,
                             folder_id INT NOT NULL,
                             user_id INT NOT NULL,
                             permission ENUM('read', 'write', 'admin') NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                             FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE,
                             FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                             UNIQUE KEY (folder_id, user_id)
);

CREATE TABLE IF NOT EXISTS files (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL,
                                     path VARCHAR(255) UNIQUE NOT NULL,
                                     size BIGINT NOT NULL,
                                     folder_id INT NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                     FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS files (
                                     id INT AUTO_INCREMENT PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL,
                                     path VARCHAR(255) UNIQUE NOT NULL,
                                     size BIGINT NOT NULL,
                                     folder_id INT NOT NULL,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                     FOREIGN KEY (folder_id) REFERENCES folders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS file_shares (
                                           id INT AUTO_INCREMENT PRIMARY KEY,
                                           file_id INT NOT NULL,
                                           share_link VARCHAR(255) UNIQUE NOT NULL,
                                           access_type ENUM('read', 'write') DEFAULT 'read',
                                           expiration DATETIME DEFAULT NULL,
                                           password VARCHAR(255) DEFAULT NULL,
                                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                                           FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);
