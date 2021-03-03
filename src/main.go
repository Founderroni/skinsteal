package main

import (
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"image"
	"image/png"
	"os"
	"strings"
)

var (
	ip   string
	port string
)

// ty TwistedAsylum in the gophertunnel discord
func SkinToRGBA(s protocol.Skin) *image.RGBA {
	t := image.NewRGBA(image.Rect(0, 0, int(s.SkinImageHeight), int(s.SkinImageWidth)))
	t.Pix = s.SkinData
	return t
}

func CapeToRGBA(s protocol.Skin) *image.RGBA {
	t := image.NewRGBA(image.Rect(0, 0, int(s.CapeImageHeight), int(s.CapeImageWidth)))
	t.Pix = s.CapeData
	return t
}

/* some servers do this cringe thing where they use SkinData
   to display information however png.Encode fails to decode it,
   most of these bots have either no name or one with fancy characters
   hence this check, if it interferes with your usage just remove the
   if isBot(name) { return } */
func isBot(s string) bool {
	if strings.Contains(s, "§") {
		return true
	} else if s == "" {
		return true
	}
	return false
}

func main() {
	fmt.Print("Server IP: ")
	_, _ = fmt.Scanln(&ip)
	fmt.Print("Server Port: ")
	_, _ = fmt.Scanln(&port)
	_ = os.Mkdir("stolen", 0755)
	dialer := minecraft.Dialer{
		TokenSource: auth.TokenSource,
	}

	address := ip + ":" + port
	conn, err := dialer.Dial("raknet", address)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	_ = conn.DoSpawn()

	for {
		pk, err := conn.ReadPacket()
		if err != nil {
			break
		}

		switch p := pk.(type) {
		case *packet.PlayerList:
			go func() {
				for _, player := range p.Entries {
					name := player.Username
					if isBot(name) {
						return
					}
					skin := SkinToRGBA(player.Skin)
					cape := CapeToRGBA(player.Skin)
					path, _ := os.Getwd()
					_ = os.Mkdir(fmt.Sprintf("%s\\stolen\\%s", path, name), 0755)
					fileSkin, _ := os.Create(fmt.Sprintf("%s\\stolen\\%s\\skin.png", path, name))
					fileCape, _ := os.Create(fmt.Sprintf("%s\\stolen\\%s\\cape.png", path, name))
					_ = png.Encode(fileSkin, skin)
					_ = png.Encode(fileCape, cape)
					fileSkin.Close()
					fileCape.Close()
					fmt.Println("Stolen " + name)
				}
			}()
		}

		p := &packet.RequestChunkRadius{ChunkRadius: 32}
		if err := conn.WritePacket(p); err != nil {
			break
		}
	}
}
