CREATE TABLE CotionUser
(
  UserID    serial                             			NOT NULL PRIMARY KEY,
  Username  varchar(20)      DEFAULT 'Username'         NOT NULL,
  Email     varchar(20)      DEFAULT 'email@email.com'  NOT NULL,
  Password  varchar(256)     DEFAULT 'password'         NOT NULL
);

CREATE TABLE Note
(
  NoteID    serial                                      NOT NULL PRIMARY KEY,
  Name      varchar(100)       DEFAULT 'Notename',
  Body      varchar(1000000) DEFAULT 'Notebody'
);

CREATE TABLE UsersNotes
(
  UserID      int        REFERENCES CotionUser 	(UserID) ON UPDATE CASCADE ON DELETE CASCADE,
  NoteID      int        REFERENCES Note 		(NoteID) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT  UserNoteID PRIMARY KEY (UserID, NoteID)
);
