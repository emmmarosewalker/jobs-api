USE jobs;
SET foreign_key_checks = 0;
DROP TABLE IF EXISTS listings,education,appuser;

CREATE TABLE listings (
  ID                INT AUTO_INCREMENT NOT NULL,
  company           VARCHAR(128) NOT NULL,
  title             VARCHAR(255) NOT NULL,
  job_type          VARCHAR(64) NOT NULL,
  job_description   TEXT NOT NULL,
  category          VARCHAR(255),
  street_address    VARCHAR(255),
  city              VARCHAR(128),
  country           VARCHAR(128),
  begin_date        DATE,
  compensation      VARCHAR(24),
  deleted_at		DATE,
  created_at		DATE,
  updated_at		DATE,

  PRIMARY KEY (`ID`)
);

INSERT INTO listings
  (company, title, job_type, job_description, category, street_address, city, country, begin_date, compensation, deleted_at)
VALUES
  ('Apple Inc', 'iOS Developer', 'Full Time', 'Were looking for an iOS Engineer to build on next generation identity product and platform. Since we are building and improving our truly groundbreaking products it would be helpful to have a strong interest and curiosity in one of the following areas: Identification, Authentication, Security and Cryptography. Experience with low level programming language is also a plus and will be useful in developing additional security features of our products. Collaborator with strong communication skills will play a significant part in working multi-functionally with other teams at Apple and external partners.', 'Software', '123 Fake St', 'Sydney', 'Australia', '2022-03-03', "175,000 / year", null);

CREATE TABLE education (
	ID 				INT AUTO_INCREMENT NOT NULL,
	degree			VARCHAR(128),
	field			VARCHAR(128),
	university		VARCHAR(128),
	city			VARCHAR(128),
	country			VARCHAR(128),
	begin_date		DATE,
	end_date		DATE,

	PRIMARY KEY (`ID`)
);

INSERT INTO education
	(degree, field, university, city, country, begin_date, end_date)
VALUES
	('Bachelor of Science', 'Computer Science', 'Macquarie University', 'Sydney', 'Australia', '2016-06-01', '2019-06-01');

CREATE TABLE appuser (
	ID				INT AUTO_INCREMENT NOT NULL,
	first_name		VARCHAR(128) NOT NULL,
	last_name		VARCHAR(128) NOT NULL,
	email			VARCHAR(128) NOT NULL,
	city			VARCHAR(128),
	country			VARCHAR(128),
	phone_number	VARCHAR(128),
	education		INT,

	PRIMARY KEY (`ID`),
	FOREIGN KEY (`education`) REFERENCES education(ID) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO appuser
	(first_name, last_name, email, city, country, phone_number, education)
VALUES
	('Emma', 'Walker', 'emma@example.com', 'Sydney', 'Australia', '+61400111222', (SELECT ID from education where country='Australia'));
