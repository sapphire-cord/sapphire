package sapphire

import (
  "github.com/bwmarrin/discordgo"
  "fmt"
  "time"
)

// Emoji constants.
const (
  EmojiLeft = "◀" // Go left, -1 page.
  EmojiRight = "▶" // Go right, +1 page.
  EmojiFirst = "⏪" // Go to first page.
  EmojiLast = "⏩" // Go to last page.
  EmojiStop = "⏹" // Stop the paginator.
)

type Paginator struct {
  Running bool // If we are running or not.
  Session *discordgo.Session // The discordgo session.
  ChannelID string // The ID of the channel we are on.
  Template func() *Embed // Base template that is passed to AddPage calls.
  Pages []*discordgo.MessageEmbed // Embeds for all pages.
  Index int // Index of current page
  Message *discordgo.Message // The sent message to be edited as we go
  AuthorID string // The user that can control this paginator.
  StopChan chan bool // Stop paginator by sending to this channel.
  Timeout time.Duration // Duration of when the paginator expires. (default: 5minutes)
}

// NewPaginator creates a new paginator and returns it.
// This is the raw one if you have special needs, it's preferred to use NewPaginatorForContext
// Session is the discordgo sesssion, channel is the ID of the channel to start the paginator.
// Author is the id of the user to listen to everyone else's reaction is ignored, pass "" to allow everyone.
func NewPaginator(session *discordgo.Session, channel, author string) *Paginator {
  return &Paginator{
    Session: session,
    ChannelID: channel,
    Running: false,
    Index: 0,
    Message: nil,
    AuthorID: author,
    StopChan: make(chan bool),
    Timeout: time.Minute * 5,
    Template: func() *Embed { return NewEmbed() },
  }
}

// NewPaginatorForContext creates a new paginator for this command context
func NewPaginatorForContext(ctx *CommandContext) *Paginator {
  return NewPaginator(ctx.Session, ctx.Channel.ID, ctx.Author.ID)
}

// SetTemplate sets the base template.
func (p *Paginator) SetTemplate(em func() *Embed) {
  p.Template = em
}

// Adds a page, takes a function that recieves the copy of embed template
// inside you can modify the embed as needed then return it back.
func (p *Paginator) AddPage(fn func(em *Embed) *Embed) {
  p.Pages = append(p.Pages, fn(p.Template()).Build())
}

// Adds a page as string, this calls the regular AddPage with the callback
// as a simple function that only sets the description to the given string.
func (p *Paginator) AddPageString(str string) {
  p.AddPage(func(em *Embed) *Embed {
    return em.SetDescription(str)
  })
}

// Add all the reactions in order.
// Called by Run to initialize.
func (p *Paginator) addReactions() {
  if p.Message == nil { return }
  p.Session.MessageReactionAdd(p.ChannelID, p.Message.ID, EmojiFirst)
  p.Session.MessageReactionAdd(p.ChannelID, p.Message.ID, EmojiLeft)
  p.Session.MessageReactionAdd(p.ChannelID, p.Message.ID, EmojiStop)
  p.Session.MessageReactionAdd(p.ChannelID, p.Message.ID, EmojiRight)
  p.Session.MessageReactionAdd(p.ChannelID, p.Message.ID, EmojiLast)
}

// Stops the paginator by sending the signal to the Stop Channel.
func (p *Paginator) Stop() {
  p.StopChan <- true
}

// Retrieves the next index for the next page
// returns 0 to go back to first page if we are on last page already.
func (p *Paginator) getNextIndex() int {
  if p.Index >= len(p.Pages) - 1 {
    return 0
  }
  return p.Index + 1
}

// Retrieves the previous index for the previous page 
// returns the last page if we are already on the first page.
func (p *Paginator) getPreviousIndex() int {
  if p.Index == 0 {
    return len(p.Pages) - 1
  }
  return p.Index - 1
}

// Sets the footers of all pages to their page number out of total pages.
// Called by Run to initialize.
func (p *Paginator) SetFooter() {
  for index, embed := range p.Pages {
    embed.Footer = &discordgo.MessageEmbedFooter{
      Text: fmt.Sprintf("Page %d/%d", index + 1, len(p.Pages)),
    }
  }
}

// Switches pages, index is assumed to be a valid index. (can panic if it's not)
// Edits the current message to the given page and updates p.Index
func (p *Paginator) Goto(index int) {
  page := p.Pages[index]
  p.Session.ChannelMessageEditEmbed(p.ChannelID, p.Message.ID, page)
  p.Index = index
}

// Switches to next page, this is safer than raw Goto as it compares indices
// and switch to first page if we are already on last one.
func (p *Paginator) NextPage() {
  p.Goto(p.getNextIndex())
}

// Switches to the previous page, this is safer than raw Goto as it compares indices
// and switch to last page if we are already on the first one.
func (p *Paginator) PreviousPage() {
  p.Goto(p.getPreviousIndex())
}

func (p *Paginator) nextReaction() chan *discordgo.MessageReactionAdd {
  channel := make(chan *discordgo.MessageReactionAdd)
  p.Session.AddHandlerOnce(func(_ *discordgo.Session, r *discordgo.MessageReactionAdd) {
    channel <- r
  })
  return channel
}

func (p *Paginator) Run() {
  if p.Running { return }
  if len(p.Pages) == 0 { return }
  p.SetFooter()
  msg, err := p.Session.ChannelMessageSendEmbed(p.ChannelID, p.Pages[0])
  if err != nil { return }
  p.Message = msg
  p.addReactions()
  p.Running = true
  start := time.Now()
  var r *discordgo.MessageReaction

  defer func() {
    p.Running = false
  }()

  for {
    select {
    case e := <-p.nextReaction():
      r = e.MessageReaction
    case <-time.After(start.Add(p.Timeout).Sub(time.Now())):
      p.Session.MessageReactionsRemoveAll(p.ChannelID, p.Message.ID)
      return
    case <-p.StopChan:
      return
    }

    if r.MessageID != p.Message.ID {
      continue
    }
    if p.AuthorID != "" && r.UserID != p.AuthorID {
      continue
    }

    go func() {
      switch r.Emoji.Name {
      case EmojiStop:
        p.Stop()
        p.Session.MessageReactionsRemoveAll(p.ChannelID, p.Message.ID)
      case EmojiRight:
        p.NextPage()
      case EmojiLeft:
        p.PreviousPage()
      case EmojiFirst:
        p.Goto(0)
      case EmojiLast:
        p.Goto(len(p.Pages) - 1)
      }
    }()
    go func() {
      time.Sleep(time.Millisecond * 250)
      p.Session.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.Name, r.UserID)
    }()
  }
}
