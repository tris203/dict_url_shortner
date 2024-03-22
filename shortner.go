package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Dictionary []string

type ShortenedChunk = struct {
	dictIdx   string
	dictVal   string
	stringIdx int
	stringVal string
}
type ShortenedChunks []ShortenedChunk

var commonWords = []string{
	"http://",
	"https://",
	"www.",
	".com",
}

const tlds = ".com,.org,.net,.int,.edu,.gov,.mil,.arpa,.ac,.ad,.ae,.af,.ag,.ai,.al,.am,.an,.ao,.aq,.ar,.as,.at,.au,.aw,.ax,.az,.ba,.bb,.bd,.be,.bf,.bg,.bh,.bi,.bj,.bl,.bm,.bn,.bo,.bq,.br,.bs,.bt,.bv,.bw,.by,.bz,.ca,.cc,.cd,.cf,.cg,.ch,.ci,.ck,.cl,.cm,.cn,.co,.cr,.cu,.cv,.cw,.cx,.cy,.cz,.de,.dj,.dk,.dm,.do,.dz,.ec,.ee,.eg,.eh,.er,.es,.et,.eu,.fi,.fj,.fk,.fm,.fo,.fr,.ga,.gb,.gd,.ge,.gf,.gg,.gh,.gi,.gl,.gm,.gn,.gp,.gq,.gr,.gs,.gt,.gu,.gw,.gy,.hk,.hm,.hn,.hr,.ht,.hu,.id,.ie,.il,.im,.in,.io,.iq,.ir,.is,.it,.je,.jm,.jo,.jp,.ke,.kg,.kh,.ki,.km,.kn,.kp,.kr,.kw,.ky,.kz,.la,.lb,.lc,.li,.lk,.lr,.ls,.lt,.lu,.lv,.ly,.ma,.mc,.md,.me,.mf,.mg,.mh,.mk,.ml,.mm,.mn,.mo,.mp,.mq,.mr,.ms,.mt,.mu,.mv,.mw,.mx,.my,.mz,.na,.nc,.ne,.nf,.ng,.ni,.nl,.no,.np,.nr,.nu,.nz,.om,.pa,.pe,.pf,.pg,.ph,.pk,.pl,.pm,.pn,.pr,.ps,.pt,.pw,.py,.qa,.re,.ro,.rs,.ru,.rw,.sa,.sb,.sc,.sd,.se,.sg,.sh,.si,.sj,.sk,.sl,.sm,.sn,.so,.sr,.ss,.st,.su,.sv,.sx,.sy,.sz,.tc,.td,.tf,.tg,.th,.tj,.tk,.tl,.tm,.tn,.to,.tp,.tr,.tt,.tv,.tw,.tz,.ua,.ug,.uk,.um,.us,.uy,.uz,.va,.vc,.ve,.vg,.vi,.vn,.vu,.wf,.ws,.ye,.yt,.za,.zm,.zw,.հայ,.বাংলা,.бел,.бг[41],.中国,.中國,.ею,.გე,.ελ[41],.香港,.भारत,.భారత్,.ભારત,.ਭਾਰਤ,.இந்தியா,.ভারত,.ಭಾರತ,.ഭാരതം,.ভাৰত,.ଭାରତ,.भारतम्,.भारोत,.қаз,.澳门,.澳門,.мкд,.мон,.рф,.срб,.新加坡,.சிங்கப்பூர்,.한국,.ලංකා,.இலங்கை,.台湾,.台灣,.ไทย,.укр,.abc,.academy,.accountant,.accountants,.active,.actor,.ads,.adult,.aero,.agency,.airforce,.analytics,.apartments,.app,.archi,.army,.art,.associates,.attorney,.auction,.audible,.audio,.author,.auto,.autos,.aws,.baby,.bananarepublic,.band,.bank,.bar,.barefoot,.bargains,.baseball,.basketball,.beauty,.beer,.best,.bestbuy,.bet,.bible,.bid,.bike,.bingo,.bio,.biz,.black,.blackfriday,.blockbuster,.blog,.blue,.boo,.book,.boots,.bot,.boutique,.box,.broadway,.broker,.build,.builders,.business,.buy,.buzz,.cab,.cafe,.call,.cam,.camera,.camp,.cancerresearch,.capital,.car,.cards,.care,.career,.careers,.cars,.case,.cash,.casino,.catering,.catholic,.center,.ceo,.cfd,.channel,.chat,.cheap,.christmas,.church,.cipriani,.circle,.city,.claims,.cleaning,.click,.clinic,.clothing,.cloud,.club,.coach,.codes,.coffee,.college,.community,.company,.compare,.computer,.condos,.construction,.consulting,.contact,.contractors,.cooking,.cool,.coop,.country,.coupon,.coupons,.courses,.credit,.creditcard,.cruise,.cricket,.cruises,.dad,.dance,.data,.date,.dating,.day,.deal,.deals,.degree,.delivery,.democrat,.dental,.dentist,.design,.dev,.diamonds,.diet,.digital,.direct,.directory,.discount,.diy,.docs,.doctor,.dog,.domains,.dot,.download,.drive,.duck,.earth,.eat,.eco,.education,.email,.energy,.engineer,.engineering,.enterprises,.equipment,.esq,.estate,.events,.exchange,.expert,.exposed,.express,.fail,.faith,.family,.fan,.fans,.farm,.fashion,.fast,.feedback,.film,.final,.finance,.financial,.fire,.fish,.fishing,.fit,.fitness,.flights,.florist,.flowers,.fly,.foo,.food,.foodnetwork,.football,.forsale,.forum,.foundation,.free,.frontdoor,.fun,.fund,.furniture,.fyi,.gallery,.game,.games,.garden,.gift,.gifts,.gives,.glass,.global,.gold,.golf,.gop,.graphics,.green,.gripe,.grocery,.group,.guide,.guitars,.guru,.hair,.hangout,.health,.healthcare,.help,.here,.hiphop,.hiv,.hockey,.holdings,.holiday,.homegoods,.homes,.homesense,.horse,.hospital,.host,.hosting,.hot,.hotels,.house,.how,.ice,.industries,.info,.ing,.ink,.institute[71],.insurance,.insure,.international,.investments,.jewelry,.jobs,.joy,.kim,.kitchen,.land,.latino,.law,.lawyer,.lease,.legal,.lgbt,.life,.lifeinsurance,.lighting,.like,.limited,.limo,.link,.live,.living,.loan,.loans,.locker,.lol,.lotto,.love,.luxury,.makeup,.management,.map,.market,.marketing,.markets,.mba,.med,.media,.meet,.meme,.memorial,.men,.menu,.mint,.mobi,.mobile,.mobily,.moe,.mom,.money,.mortgage,.motorcycles,.mov,.movie,.museum,.name,.navy,.network,.new,.news,.ngo,.ninja,.now,.observer,.off,.one,.ong,.onl,.online,.ooo,.open,.organic,.origins,.page,.partners,.parts,.party,.pay,.pet,.pharmacy,.phone,.photo,.photography,.photos,.physio,.pics,.pictures,.pid,.pin,.pink,.pizza,.place,.plumbing,.plus,.poker,.porn,.post,.press,.prime,.pro,.productions,.prof,.promo,.properties,.property,.protection,.pub,.qpon,.racing,.radio,.read,.realestate,.realtor,.realty,.recipes,.red,.rehab,.reit,.ren,.rent,.rentals,.repair,.report,.republican,.rest,.restaurant,.review,.reviews,.rich,.rip,.rocks,.rodeo,.room,.rugby,.run,.safe,.sale,.save,.scholarships,.school,.science,.search,.secure,.security,.select,.services,.sex,.sexy,.shoes,.shop,.shopping,.show,.showtime,.silk,.singles,.site,.ski,.skin,.sky,.sling,.smile,.soccer,.social,.software,.solar,.solutions,.song,.space,.spot,.spreadbetting,.storage,.store,.stream,.studio,.study,.style,.sucks,.supplies,.supply,.support,.surf,.surgery,.systems,.talk,.tattoo,.tax,.taxi,.team,.tech,.technology,.tel,.tennis,.theater,.theatre,.tickets,.tips,.tires,.today,.tools,.top,.tours,.town,.toys,.trade,.trading,.training,.travel,.travelersinsurance,.trust,.tube,.tunes,.uconnect,.university,.vacations,.ventures,.vet,.video,.villas,.vip,.vision,.vodka,.vote,.voting,.voyage,.wang,.watch,.watches,.weather,.webcam,.website,.wed,.wedding,.whoswho,.wiki,.win,.wine,.winners,.work,.works,.world,.wow,.wtf,.xxx,.xyz,.yachts,.yoga,.you,.zero,.zone,.shouji,.tushu,.wanggou,.weibo,.xihuan,.arte,.clinique,.luxe,.maison,.moi,.rsvp,.sarl,.epost,.gmbh,.haus,.immobilien,.jetzt,.kaufen,.kinder,.reise,.reisen,.schule,.versicherung,.desi,.shiksha,.casa,.immo,.moda,.voto,.bom,.passagens,.abogado,.gratis,.futbol,.hoteles,.juegos,.ltda,.soy,.tienda,.uno,.viajes,.vuelos,.كوم,.موبايلي,.كاثوليك,.بيتك,.在线,.中文网,.移动,.网址,.网络,.公司,.商城,.机构,.我爱你,.商标,.世界,.集团,.дети,.католик,.ком,.онлайн,.орг,.сайт,.संगठन,.कॉम,.नेट,.닷컴,.닷넷,.קום‎,.みんな,.セール,.ファッション,.ストア,.ポイント,.クラウド,.コム,.คอม,.africa,.capetown,.durban,.joburg,.abudhabi,.arab,.asia,.doha,.dubai,.krd,.kyoto,.nagoya,.okinawa,.osaka,.ryukyu,.taipei,.tatar,.tokyo,.yokohama,.alsace,.amsterdam,.bcn,.barcelona,.bayern,.berlin,.brussels,.budapest,.bzh,.cat,.cologne,.corsica,.cymru,.eus,.frl,.gal,.gent,.hamburg,.helsinki,.irish,.ist,.istanbul,.koeln,.london,.madrid,.moscow (ru),.nrw,.paris,.ruhr,.saarland,.scot,.stockholm,.swiss,.tirol,.vlaanderen,.wales,.wien,.zuerich,.boston,.miami,.nyc,.quebec,.vegas,.kiwi,.melbourne,.sydney,.lat,.rio,.佛山,.广东,.москва (ru),.moscow,.рус (ru),.ابوظبي,.aaa,.aarp,.abarth,.abb,.abbott,.abbvie,.accenture,.aco,.aeg,.aetna,.afl,.agakhan,.aig,.aigo,.airbus,.airtel,.akdn,.alfaromeo,.alibaba,.alipay,.allfinanz,.allstate,.ally,.alstom,.americanexpress,.amex,.amica,.android,.anz,.aol,.apple,.aquarelle,.aramco,.audi,.auspost,.axa,.azure,.baidu,.barclaycard,.barclays,.bauhaus,.bbc,.bbt,.bbva,.bcg,.bentley,.bharti,.bing,.blanco,.bloomberg,.bms,.bmw,.bnl,.bnpparibas,.boehringer,.bond,.booking,.bosch,.bostik,.bradesco,.bridgestone,.brother,.bugatti,.cal,.calvinklein,.canon,.capitalone,.caravan,.cartier,.cba,.cbn,.cbre,.cbs,.cern,.cfa,.chanel,.chase,.chintai,.chrome,.chrysler,.cisco,.citadel,.citi,.citic,.clubmed,.comcast,.commbank,.creditunion,.crown,.crs,.cuisinella,.dabur,.datsun,.dealer,.dell,.deloitte,.delta,.dhl,.discover,.dish,.dnp,.dodge,.dunlop,.dupont,.dvag,.edeka,.emerck,.epson,.ericsson,.erni,.esurance,.etisalat,.eurovision,.everbank,.extraspace,.fage,.fairwinds,.farmers,.FedEx,.ferrari,.ferrero,.fiat,.fidelity,.firestone,.firmdale,.flickr,.flir,.flsmidth,.ford,.fox,.fresenius,.forex,.frogans,.frontier,.fujitsu,.fujixerox,.gallo,.gallup,.gap,.gbiz,.gea,.genting,.giving,.gle,.globo,.gmail,.gmo,.gmx,.godaddy,.goldpoint,.goodyear,.google,.grainger,.guardian,.gucci,.hbo,.hdfc,.hdfcbank,.hermes,.hisamitsu,.hitachi,.hkt,.honda,.honeywell,.hotmail,.hsbc,.htc,.hughes,.hyatt,.hyundai,.ibm,.ieee,.ifm,.ikano,.imdb,.infiniti,.intel,.intuit,.ipiranga,.iselect,.itau,.itv,.iveco,.jaguar,.java,.jcb,.jcp,.jeep,.jpmorgan,.juniper,.kddi,.kerryhotels,.kerrylogistics,.kerryproperties,.kfh,.kia,.kindle,.komatsu,.kpmg,.kred,.kuokgroup,.lacaixa,.ladbrokes,.lamborghini,.lancaster,.lancia,.lancome,.landrover,.lanxess,.lasalle,.latrobe,.lds,.lego,.liaison,.lexus,.lidl,.lifestyle,.lilly,.lincoln,.linde,.lipsy,.lixil,.locus,.lotte,.lpl,.lplfinancial,.lundbeck,.lupin,.macys,.maif,.man,.mango,.marriott,.maserati,.mattel,.mcd,.mcdonalds,.mckinsey,.meo,.metlife,.microsoft,.mini,.mit,.mitsubishi,.mlb,.mma,.monash,.mormon,.moto,.movistar,.msd,.mtn,.mtpc,.mtr,.mutual,.mutuelle,.nadex,.nationwide,.natura,.nba,.nec,.netflix,.neustar,.newholland,.nexus,.nfl,.nhk,.nico,.nike,.nikon,.nissan,.nissay,.nokia,.northwesternmutual,.norton,.nra,.ntt,.obi,.office,.omega,.oracle,.orange,.orientexpress,.otsuka,.ovh,.pamperedchef,.panasonic,.pccw,.pfizer,.philips,.piaget,.pictet,.ping,.pioneer,.play,.playstation,.pohl,.politie,.praxi,.prod,.progressive,.pru,.prudential,.pwc,.quest,.qvc,.redstone,.reliance,.rexroth,.ricoh,.rmit,.rocher,.rogers,.rwe,.safety,.sakura,.samsung,.sandvik,.sandvikcoromant,.sanofi,.sap,.saxo,.sbi,.sbs,.sca,.scb,.schaeffler,.schmidt,.schwarz,.scjohnson,.scor,.seat,.sener,.ses,.sew,.seven,.sfr,.seek,.shangrila,.sharp,.shaw,.shell,.shriram,.sina,.skype,.smart,.sncf,.softbank,.sohu,.Sony,.spiegel,.stada,.staples,.star,.starhub,.statebank,.statefarm,.statoil,.stc,.stcgroup,.suzuki,.swatch,.swiftcover,.symantec,.taobao,.target,.tatamotors,.tdk,.telecity,.telefonica,.temasek,.teva,.tiffany,.tjx,.toray,.toshiba,.total,.toyota,.travelchannel,.travelers,.tui,.tvs,.ubs,.unicom,.uol,.ups,.vanguard,.verisign,.vig,.viking,.virgin,.visa,.vista,.vistaprint,.vivo,.volkswagen,.volvo,.walmart,.walter,.weatherchannel,.weber,.weir,.williamhill,.windows,.wme,.wolterskluwer,.woodside,.wtc,.xbox,.xerox,.xfinity,.xperia,.yahoo,.yamaxun,.yandex,.yodobashi,.youtube,.zappos,.zara,.zip,.zippo,.ارامكو,.联通,.中信,.香格里拉,.淡马锡,.大众汽车,.vermögensberater,.vermögensberatung,.グーグル,.谷歌,.工行,.嘉里,.嘉里大酒店,.飞利浦,.诺基亚,.電訊盈科,.삼성,.example,.invalid,.local,.localhost,.onion,.test,.i2p,.glue,.测试,.測試,.испытание,.परीक्षा,.δοκιμή,.테스트,.テスト,.பரிட்சை'"
const extensions = ".html,.asp,apsx,.php,.jsp,.htm,.txt"
const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"


func getDictionary() Dictionary {
	jsonFile, err := os.Open("dictionary.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened dictionary.json")

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var words Dictionary
	json.Unmarshal(byteValue, &words)
	fmt.Println("Loaded ", len(words), " words")

	sort.Slice(words, func(i, j int) bool {
		return len(words[i]) > len(words[j])
	})
	var addWords = []string{}

	addWords = append(strings.Split(tlds, ","), addWords...)
	addWords = append(strings.Split(extensions, ","), addWords...)
	addWords = append(commonWords, addWords...)

	// padd addWord to 10,000 to presrve index
	for i := len(addWords); i < 10000; i++ {
		addWords = append(addWords, "YOUWILLNEVERFINDMEINAURLORSTUFFAOIJDOJOIAWJDOIJDWIJDWIJ")
	}

	words = append(words, addWords...)

	// append charset to the dictionary
	words = append(words, strings.Split(charset, "")...)

	// add common words to the dictionary at the beginning
	fmt.Println("Adding common words to the dictionary", len(addWords))
	fmt.Println("Total words in dictionary: ", len(words))

	return words
}

func shortenURL(url string, words Dictionary) string {

	var replacedCharCount int
	var chunks ShortenedChunks
	for replacedCharCount < len(url) {
		var remainingString = url[replacedCharCount:]
		var found = false
		for i := 0; i < len(words); i++ {
			if strings.HasPrefix(remainingString, words[i]) {
				chunks = append(chunks, ShortenedChunk{
					dictIdx:   base62Encode(i),
					dictVal:   words[i],
					stringIdx: replacedCharCount,
					stringVal: url[replacedCharCount : replacedCharCount+len(words[i])],
				})
				replacedCharCount += len(words[i])
				found = true
				break
			}

		}
		if found == false {
			//if no match found, add the character to the chunks
			chunks = append(chunks, ShortenedChunk{
				dictIdx:   "0",
				dictVal:   url[replacedCharCount : replacedCharCount+1],
				stringIdx: replacedCharCount,
				stringVal: url[replacedCharCount : replacedCharCount+1],
			})
			replacedCharCount++
			continue
		}

	}

	var result string
	for i := 0; i < len(chunks); i++ {
		if chunks[i].dictIdx == "0" {
			result += chunks[i].stringVal
		} else {
			result += chunks[i].dictIdx
		}
	}
	return result

}

func expandUrl(shortened string, words Dictionary) string {
	var result string
	for i := 0; i < len(shortened); i += 3 {
		// if the chunk contains a character that is not in the charset, add it to the result
		if strings.Index(charset, string(shortened[i])) == -1 {
			result += string(shortened[i])
			i -= 2
			continue
		}
		var chunk = shortened[i : i+3]
		var idx = base62Decode(chunk)
		result += words[idx]
	}
	return result
}

func base62Encode(n int) string {
	var result string
	for n > 0 {
		result = string(charset[n%62]) + result
		n = n / 62
	}
	//pad to 3 characters
	for len(result) < 3 {
		result = "0" + result
	}
	return result
}

func base62Decode(s string) int {
	var result int
	for i := 0; i < len(s); i++ {
		result = result*62 + strings.Index(charset, string(s[i]))
	}
	return result
}
