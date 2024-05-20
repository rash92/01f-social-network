import { useEffect, useRef,  } from "react"

export default function Websocket(){
    // const [isPaused, setPause] = useState(false)
    const ws = useRef(null)

    useEffect(()=>{
        ws.current = new WebSocket("ws://localhost:8080/ws")

        ws.current.onopen = () => console.log("ws opened")
        ws.current.onclose = () => console.log("ws closed")
        ws.current.onmessage = e => {
            console.log("message: ", e.data)
        }
        const wsCurrent = ws.current

        return () => wsCurrent.close()
    },[])

    return (
    <div>
        Websocket here!
    </div>
    )
}