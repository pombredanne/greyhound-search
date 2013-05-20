/*
Code for a greyhound search server

TBD
*/
package greyhound

import "encoding/json"
import "fmt"
import "log"
import "net/http"
import "code.google.com/p/go.net/websocket"

func (gs *GreyhoundSearch) HandleGreyhoundSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	action, hasAction := r.Form["action"]
	if !hasAction {
		fmt.Fprintf(w, "no action argument passed!")
	} else {
		queryData := make(map[string]string)
		for k, v := range r.Form {
			queryData[k] = v[0]
		}
		msg := &Message{action[0], queryData}
		fmt.Fprintf(w, gs.PerformAction(msg))
	}
}

/* handles greyhound-search's websocket actions.
effectively, greyhound messages are always sent as json. Specifically:
{ action: 'ACTION',
  data: { JSON_OBJECT }
}

each action has a struct to unmarshal json, and returns a series of values
 */
func (gs *GreyhoundSearch) HandleGreyhoundSearchSocket(ws *websocket.Conn) {
	for {
		var msg Message
		err := websocket.JSON.Receive(ws, &msg)
		log.Println("raw message: ",  msg)
		if err != nil { 
			fmt.Println(err)
			break
		}
		_ = websocket.Message.Send(ws, gs.PerformAction(&msg))
	}
}

func (gs *GreyhoundSearch) PerformAction (m *Message) string {
	var out_json []byte
	switch m.Action {
	case "query": 
		out_json, _ = json.Marshal(gs.Search(m.QueryData["project"], m.QueryData["query"]))
  case "list_projects":
		out_json, _ = json.Marshal(gs.ListProjects())
	case "view_file":
		out_json, _ = json.Marshal(gs.ViewFile(m.QueryData["file"]))
	default:
		out_json, _ = json.Marshal([]string{fmt.Sprintf("%s is not a valid action", m.Action)})
	}
	return string(out_json)
}

type Message struct {
	Action string
	QueryData map[string]string
}
