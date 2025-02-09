DROP DATABASE IF EXISTS cnadscdb;
CREATE DATABASE cnadscdb;
USE cnadscdb;
DROP TABLE IF EXISTS Senior ;
DROP TABLE IF EXISTS Emergency_Contact ;
DROP TABLE IF EXISTS Assessment ;
DROP TABLE IF EXISTS Emergency_Contact ;
DROP TABLE IF EXISTS HealthGuide;

-- Senior Table
CREATE TABLE Senior (
    SeniorID  INT PRIMARY KEY AUTO_INCREMENT,
    Phone_number VARCHAR(15) NOT NULL,
    Name VARCHAR(100)
);

CREATE TABLE Emergency_Contact(
Contactid INT PRIMARY KEY AUTO_INCREMENT,
ContactName VARCHAR(255),
ContactNumber INT NOT NULL,
SeniorID INT NOT NULL,
FOREIGN KEY (SeniorID) REFERENCES Senior(SeniorID)
);
 
-- Assessment Table
CREATE TABLE Assessment (
AssessmentID INT PRIMARY KEY AUTO_INCREMENT,
Overall_Wellbeing VARCHAR(255) NOT NULL,
    SeniorID INT NOT NULL, FOREIGN KEY (SeniorID) REFERENCES Senior(SeniorID)
); 

-- Clinic Table
CREATE TABLE Clinic (
    ClinicID INT PRIMARY KEY AUTO_INCREMENT,
    Name VARCHAR(100) NOT NULL,
    Location VARCHAR(255) NOT NULL,
    Contact INT NOT NULL
);
 
-- Doctor Table
CREATE TABLE Doctor (
    DoctorID INT PRIMARY KEY AUTO_INCREMENT,
    Name VARCHAR(255) NOT NULL,
    Contact VARCHAR(15) NOT NULL,
    ClinicID INT NOT NULL,
    Specialization VARCHAR(100) NOT NULL,
    Email VARCHAR(100) UNIQUE NOT NULL,
    Password VARCHAR(255) NOT NULL,
    FOREIGN KEY (ClinicID) REFERENCES Clinic(ClinicID)
);
 
-- Conditions Table
CREATE TABLE Conditions (
    ConditionID INT PRIMARY KEY AUTO_INCREMENT,
    Condition_Name VARCHAR(100) NOT NULL,
    Condition_Description TEXT NOT NULL
);
 
-- Health Condition Table
CREATE TABLE HealthCondition (
    HealthconditionID INT PRIMARY KEY AUTO_INCREMENT,
    ConditionID INT NOT NULL,
    SeniorID INT NOT NULL,
    Onset_datetime DATETIME NOT NULL,
    FOREIGN KEY (ConditionID) REFERENCES Conditions(ConditionID),
    FOREIGN KEY (SeniorID) REFERENCES Senior(SeniorID)
);
 
-- Suggestion Table
CREATE TABLE Suggestion (
    SuggestionID INT PRIMARY KEY AUTO_INCREMENT,
    ConditionID INT NOT NULL,
    Suggestion TEXT NOT NULL,
    SeniorID INT,
    FOREIGN KEY (ConditionID) REFERENCES Conditions(ConditionID),
    FOREIGN KEY (SeniorID) REFERENCES Senior(SeniorID)
);
 
-- Call Record Table
CREATE TABLE CallRecord (
    CallRecordID INT PRIMARY KEY AUTO_INCREMENT,
    DoctorID INT NOT NULL,
    SeniorID INT NOT NULL,
    Call_Datetime DATETIME NOT NULL,
    Duration INT NOT NULL,
    FOREIGN KEY (DoctorID) REFERENCES Doctor(DoctorID),
    FOREIGN KEY (SeniorID) REFERENCES Senior(SeniorID)
);


-- Health Guides
CREATE TABLE HealthGuide (
    HealthGuideID INT PRIMARY KEY AUTO_INCREMENT,
    HealthGuideDescription TEXT NOT NULL,
    HealthGuideVideoLink VARCHAR(255) NOT NULL,
    Overall_Wellbeing VARCHAR(255) NOT NULL
);