CREATE TABLE CotionUser
(
  UserID    varchar(64)      NOT NULL PRIMARY KEY,
  Username  varchar(20)      NOT NULL,
  Email     varchar(20)      NOT NULL,
  Password  varchar(256)     NOT NULL
);

CREATE TABLE Note
(
  NoteID    varchar(100)     PRIMARY KEY,
  Name      varchar(100)     NOT NULL,
  Body      text             NOT NULL
);

CREATE TABLE UsersNotes
(
  UserID      varchar(64)        REFERENCES CotionUser 	(UserID) ON UPDATE CASCADE ON DELETE CASCADE,
  NoteID      varchar(100)       REFERENCES Note 		(NoteID) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT  UserNoteID PRIMARY KEY (UserID, NoteID)
);
