Tämä toimii lyhyenä selitteenä, miten File_Microservices toimii.

Perusidea:

Luodaan aluksi yksinkertainen palvelin, joka käsittelee HTTP kutsuja.

Testidata.go:ssa on data storage (tässä tapauksessa slice), joka sisältää testidataa.

Testidata on tyyppiä File, jolla on neljä kenttää ID, Nimi, Tyyppi ja Sisältö
Testitiedostot luodaan myös lokaalisti.

Tähän dataan päästään käsiksi HTTP kutsuilla POST, PUT, DELETE ja GET.
    POST kutsulla voidaan luoda uusia Filejä
    DELETE kutsulla voidaan poistaa olemassa olevia Filejä
    GET kutsulla voidaan hakea kaikki olemassa olevat Filet

Esimerkit kutsuista (win10 ympäristössä):

curl -X POST -H "Content-Type:application/json" -d "{\"name\": \"testifilu\",\"type\": \".txt\,\"content\": \"sisältöä\"}" localhost:8080

curl -X DELETE localhost:8080/1

GET on oletusarvo, joten siihen ei tarvita -X GET
curl localhost:8080

