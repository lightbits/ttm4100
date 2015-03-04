package coding

import(
    "encoding/json"
    "log"
)

type ClientPackage struct {
    Request string
    Content string
}

type ServerPackage struct {
    Timestamp string
    Sender string
    Response string
    Content string
}

func ClientPackageToNetworkPacket(pack ClientPackage) []byte {
    byteArr, err := json.Marshal(pack)
    if err != nil {log.Println(err)}
    return byteArr
}

func ServerPackagesToNetworkPacket(pack []ServerPackage) []byte {
    byteArr, err := json.Marshal(pack)
    if err != nil {log.Println(err)}
    return byteArr
}

func NetworkPacketToClientPackage(byteArr []byte) ClientPackage {
    var ClientPack ClientPackage
    err := json.Unmarshal(byteArr[:], &ClientPack)
    if(err != nil) {log.Println(err)}
    return ClientPack
}

func NetworkPacketToServerPackages(byteArr []byte) []ServerPackage {
    // log.Println(string(byteArr))
    var ServerPack []ServerPackage
    err := json.Unmarshal(byteArr[:], &ServerPack)
    if (err != nil) {log.Println(err)}
    return ServerPack
}
