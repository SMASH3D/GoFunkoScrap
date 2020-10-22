USE funkoscrap;

CREATE TABLE `licences` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `LicenceID` int(11) DEFAULT NULL,
  `Name` varchar(256) DEFAULT NULL,
  `Logo` varchar(256) DEFAULT NULL,
  `URL` varchar(256) DEFAULT NULL,
  `CrawledAt` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `licenceID_idx` (`LicenceID`)
) ENGINE=InnoDB AUTO_INCREMENT=361 DEFAULT CHARSET=latin1

CREATE TABLE `funkos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `LicenceID` int(11) DEFAULT NULL,
  `Ref` varchar(256),
  `Num` int(11) DEFAULT NULL,
  `Produced` datetime DEFAULT NULL,
  `Scale` varchar(32) DEFAULT NULL,
  `Name` varchar(64) DEFAULT NULL,
  `Edition` varchar(256) DEFAULT NULL,
  `ImgURL` varchar(256) DEFAULT NULL,
  `Price` decimal(7,2) DEFAULT NULL,
  `CrawledAt` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ref_idx` (`Ref`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1
