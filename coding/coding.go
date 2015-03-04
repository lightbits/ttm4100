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

func ServerPackageToNetworkPacket(pack ServerPackage) []byte {
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
    var ServerPack ServerPackage
    err := json.Unmarshal(byteArr[:], &ServerPack)
    // Oh boy...
    if (err != nil) {
        // It may have been a list instead
        var ServerPacks []ServerPackage
        err = json.Unmarshal(byteArr[:], &ServerPacks)
        if err != nil {
            log.Println(err)
        }
        return ServerPacks
    }
    return []ServerPackage{ServerPack}
}
