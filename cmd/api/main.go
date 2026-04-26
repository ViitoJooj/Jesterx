package main

import (
	"log"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
	"github.com/ViitoJooj/Jesterx/pkg/supabase"
	"github.com/ViitoJooj/Jesterx/pkg/validators"
)

func main() {
	dotenv.Set()
	supabase.Conn()
	validators.LoadEmbedded()

	log.Println("Running...")
}
