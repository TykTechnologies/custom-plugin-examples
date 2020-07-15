CREATE TABLE [dbo].[users]
(
  [id] [int]  NOT NULL,
  [username] [varchar](50)  NOT NULL,
  [token] [varchar](50)  NOT NULL
)
ALTER TABLE [dbo].[users] ADD CONSTRAINT PK__users__3213E83F655ADDA9 PRIMARY KEY  ([id])

-- Add 2 rows for users.
INSERT INTO users (id, username, token) VALUES 
(1,'user_a','b5b037a78522671b89a2c1b21d9b80c6'),
(2,'user_b','396ac2a33748bb784077f4112335faa0');

