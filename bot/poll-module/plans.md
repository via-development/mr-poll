# Poll Data Modal

- GuildId
- ChannelId
- MessageId
- UserId
- Type: Yes-or-no(0)/Single-choice(1)/Multiple-choice(2)/Submit-option(3)
- Question
- Options
  - Name
  - Id
  - Emoji
  - Voters
- AnonymousType: None(0)/Forever(1)/Until-end(2)
- RequiredRoles

# Poll Module Functions

- Create A Poll
  - Check Channel Permissions
  - Create Poll on Database
  - Create Poll Message on Discord
- Get Poll from Database
- Get List of Polls from Database
- Handle Poll Vote
  - Get Poll & Author Data from Database
  - Remove vote if already voted
  - Switch vote if different from last vote
  - Add vote if hasn't voted
  - Update Poll Message
  - 