/*M!999999\- enable the sandbox mode */ 
-- MariaDB dump 10.19-11.6.2-MariaDB, for Linux (x86_64)
--
-- Host: 192.168.128.3    Database: cnad-scdb
-- ------------------------------------------------------
-- Server version	11.4.3-MariaDB-ubu2404

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*M!100616 SET @OLD_NOTE_VERBOSITY=@@NOTE_VERBOSITY, NOTE_VERBOSITY=0 */;

--
-- Table structure for table `Assessment`
--

DROP TABLE IF EXISTS `Assessment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Assessment` (
  `AssessmentID` int(11) NOT NULL AUTO_INCREMENT,
  `Overall_Wellbeing` varchar(255) NOT NULL,
  `SeniorID` int(11) NOT NULL,
  PRIMARY KEY (`AssessmentID`),
  KEY `SeniorID` (`SeniorID`),
  CONSTRAINT `Assessment_ibfk_1` FOREIGN KEY (`SeniorID`) REFERENCES `Senior` (`SeniorID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Assessment`
--

LOCK TABLES `Assessment` WRITE;
/*!40000 ALTER TABLE `Assessment` DISABLE KEYS */;
/*!40000 ALTER TABLE `Assessment` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `CallRecord`
--

DROP TABLE IF EXISTS `CallRecord`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `CallRecord` (
  `CallRecordID` int(11) NOT NULL AUTO_INCREMENT,
  `DoctorID` int(11) NOT NULL,
  `SeniorID` int(11) NOT NULL,
  `Call_Datetime` datetime NOT NULL,
  `Duration` int(11) NOT NULL,
  PRIMARY KEY (`CallRecordID`),
  KEY `DoctorID` (`DoctorID`),
  KEY `SeniorID` (`SeniorID`),
  CONSTRAINT `CallRecord_ibfk_1` FOREIGN KEY (`DoctorID`) REFERENCES `Doctor` (`DoctorID`),
  CONSTRAINT `CallRecord_ibfk_2` FOREIGN KEY (`SeniorID`) REFERENCES `Senior` (`SeniorID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `CallRecord`
--

LOCK TABLES `CallRecord` WRITE;
/*!40000 ALTER TABLE `CallRecord` DISABLE KEYS */;
/*!40000 ALTER TABLE `CallRecord` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Clinic`
--

DROP TABLE IF EXISTS `Clinic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Clinic` (
  `ClinicID` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(100) NOT NULL,
  `Location` varchar(255) NOT NULL,
  `Contact` int(11) NOT NULL,
  PRIMARY KEY (`ClinicID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Clinic`
--

LOCK TABLES `Clinic` WRITE;
/*!40000 ALTER TABLE `Clinic` DISABLE KEYS */;
/*!40000 ALTER TABLE `Clinic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Conditions`
--

DROP TABLE IF EXISTS `Conditions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Conditions` (
  `ConditionID` int(11) NOT NULL AUTO_INCREMENT,
  `Condition_Name` varchar(100) NOT NULL,
  `Condition_Description` text NOT NULL,
  PRIMARY KEY (`ConditionID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Conditions`
--

LOCK TABLES `Conditions` WRITE;
/*!40000 ALTER TABLE `Conditions` DISABLE KEYS */;
/*!40000 ALTER TABLE `Conditions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Doctor`
--

DROP TABLE IF EXISTS `Doctor`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Doctor` (
  `DoctorID` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(255) NOT NULL,
  `Contact` varchar(15) NOT NULL,
  `ClinicID` int(11) NOT NULL,
  `Specialization` varchar(100) NOT NULL,
  `Email` varchar(100) NOT NULL,
  `Password` varchar(255) NOT NULL,
  PRIMARY KEY (`DoctorID`),
  UNIQUE KEY `Email` (`Email`),
  KEY `ClinicID` (`ClinicID`),
  CONSTRAINT `Doctor_ibfk_1` FOREIGN KEY (`ClinicID`) REFERENCES `Clinic` (`ClinicID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Doctor`
--

LOCK TABLES `Doctor` WRITE;
/*!40000 ALTER TABLE `Doctor` DISABLE KEYS */;
/*!40000 ALTER TABLE `Doctor` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Emergency_Contact`
--

DROP TABLE IF EXISTS `Emergency_Contact`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Emergency_Contact` (
  `Contactid` int(11) NOT NULL AUTO_INCREMENT,
  `ContactName` varchar(255) DEFAULT NULL,
  `ContactNumber` int(11) NOT NULL,
  `SeniorID` int(11) NOT NULL,
  PRIMARY KEY (`Contactid`),
  KEY `SeniorID` (`SeniorID`),
  CONSTRAINT `Emergency_Contact_ibfk_1` FOREIGN KEY (`SeniorID`) REFERENCES `Senior` (`SeniorID`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Emergency_Contact`
--

LOCK TABLES `Emergency_Contact` WRITE;
/*!40000 ALTER TABLE `Emergency_Contact` DISABLE KEYS */;
INSERT INTO `Emergency_Contact` VALUES
(5,'ghost in hotel',91542235,3),
(6,'Name',123456,3),
(7,'peter',11223345,3),
(8,'YQY',91542235,3),
(9,'YQY',91542235,3),
(10,'ghost in hotel',85178498,3),
(16,'YE QIYANG',85178498,4),
(17,'YQY2',91542235,4);
/*!40000 ALTER TABLE `Emergency_Contact` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `HealthCondition`
--

DROP TABLE IF EXISTS `HealthCondition`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `HealthCondition` (
  `HealthconditionID` int(11) NOT NULL AUTO_INCREMENT,
  `ConditionID` int(11) NOT NULL,
  `SeniorID` int(11) NOT NULL,
  `Onset_datetime` datetime NOT NULL,
  PRIMARY KEY (`HealthconditionID`),
  KEY `ConditionID` (`ConditionID`),
  KEY `SeniorID` (`SeniorID`),
  CONSTRAINT `HealthCondition_ibfk_1` FOREIGN KEY (`ConditionID`) REFERENCES `Conditions` (`ConditionID`),
  CONSTRAINT `HealthCondition_ibfk_2` FOREIGN KEY (`SeniorID`) REFERENCES `Senior` (`SeniorID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `HealthCondition`
--

LOCK TABLES `HealthCondition` WRITE;
/*!40000 ALTER TABLE `HealthCondition` DISABLE KEYS */;
/*!40000 ALTER TABLE `HealthCondition` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Senior`
--

DROP TABLE IF EXISTS `Senior`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Senior` (
  `SeniorID` int(11) NOT NULL AUTO_INCREMENT,
  `Phone_number` varchar(15) NOT NULL,
  `Name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`SeniorID`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Senior`
--

LOCK TABLES `Senior` WRITE;
/*!40000 ALTER TABLE `Senior` DISABLE KEYS */;
INSERT INTO `Senior` VALUES
(3,'91542235',NULL),
(4,'85178498',NULL);
/*!40000 ALTER TABLE `Senior` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Suggestion`
--

DROP TABLE IF EXISTS `Suggestion`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Suggestion` (
  `SuggestionID` int(11) NOT NULL AUTO_INCREMENT,
  `ConditionID` int(11) NOT NULL,
  `Suggestion` text NOT NULL,
  `SeniorID` int(11) DEFAULT NULL,
  PRIMARY KEY (`SuggestionID`),
  KEY `ConditionID` (`ConditionID`),
  KEY `SeniorID` (`SeniorID`),
  CONSTRAINT `Suggestion_ibfk_1` FOREIGN KEY (`ConditionID`) REFERENCES `Conditions` (`ConditionID`),
  CONSTRAINT `Suggestion_ibfk_2` FOREIGN KEY (`SeniorID`) REFERENCES `Senior` (`SeniorID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Suggestion`
--

LOCK TABLES `Suggestion` WRITE;
/*!40000 ALTER TABLE `Suggestion` DISABLE KEYS */;
/*!40000 ALTER TABLE `Suggestion` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

-- Dump completed on 2025-02-07  0:37:58
