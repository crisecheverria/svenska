package data

// Story is a short Swedish text with its English translation, for typing practice.
type Story struct {
	Title string
	Sv    string
	En    string
}

var Stories = []Story{
	{
		Title: "Min dag",
		Sv:    "Jag vaknar klockan sju varje morgon. Först dricker jag kaffe och äter frukost. Sedan går jag till jobbet med bussen. Jag jobbar på ett kontor i centrum. På lunchen äter jag en smörgås med mina kollegor. Efter jobbet går jag hem och lagar middag. På kvällen tittar jag på tv eller läser en bok innan jag somnar.",
		En:    "I wake up at seven every morning. First I drink coffee and eat breakfast. Then I go to work by bus. I work at an office in the city center. At lunch I eat a sandwich with my colleagues. After work I go home and cook dinner. In the evening I watch tv or read a book before I fall asleep.",
	},
	{
		Title: "På café",
		Sv:    "Jag går till ett café i centrum. Jag beställer en kopp kaffe och en kanelbulle. Det är mysigt att sitta och titta på folk. Kaffet smakar bra och bullen är nybakad. Jag läser tidningen medan jag dricker mitt kaffe. En vän ringer och vi bestämmer att träffas här nästa vecka också.",
		En:    "I go to a café in the city center. I order a cup of coffee and a cinnamon bun. It is cozy to sit and watch people. The coffee tastes good and the bun is freshly baked. I read the newspaper while I drink my coffee. A friend calls and we decide to meet here next week too.",
	},
	{
		Title: "Min familj",
		Sv:    "Jag har en stor familj. Min mamma heter Anna och min pappa heter Erik. Jag har två systrar och en bror. Vi bor i Stockholm men mina föräldrar bor i Uppsala. På söndagar äter vi middag tillsammans. Mamma lagar alltid något gott. Vi pratar och skrattar och barnen leker i trädgården.",
		En:    "I have a big family. My mother's name is Anna and my father's name is Erik. I have two sisters and one brother. We live in Stockholm but my parents live in Uppsala. On Sundays we eat dinner together. Mom always cooks something delicious. We talk and laugh and the children play in the garden.",
	},
	{
		Title: "I parken",
		Sv:    "Det är en vacker dag. Solen skiner och det är varmt. Jag tar en promenad i parken. Barnen leker och hundarna springer. Jag sätter mig på en bänk och läser en bok. En gammal man matar fåglarna bredvid mig. Det doftar blommor överallt. Jag stannar i parken hela eftermiddagen.",
		En:    "It is a beautiful day. The sun is shining and it is warm. I take a walk in the park. The children play and the dogs run. I sit on a bench and read a book. An old man feeds the birds next to me. It smells like flowers everywhere. I stay in the park the whole afternoon.",
	},
	{
		Title: "Mataffären",
		Sv:    "Jag behöver handla mat. Jag tar en korg och går in i affären. Jag köper mjölk, bröd, ost och frukt. Äpplena ser fina ut. Jag hittar också en god chokladkaka på rea. I kassan står det en lång kö men jag har inte bråttom. Jag betalar med kort och går hem med mina kassar.",
		En:    "I need to buy groceries. I take a basket and go into the store. I buy milk, bread, cheese and fruit. The apples look nice. I also find a good chocolate cake on sale. At the register there is a long queue but I am not in a hurry. I pay by card and walk home with my bags.",
	},
	{
		Title: "Midsommar",
		Sv:    "Midsommar är en viktig tradition i Sverige. Vi dansar runt midsommarstången och sjunger sånger. Vi äter sill, potatis och jordgubbar med grädde. Barnen plockar blommor och gör kransar. Det är den längsta dagen på året och solen går nästan aldrig ner. Vi firar med familj och vänner utomhus hela dagen.",
		En:    "Midsummer is an important tradition in Sweden. We dance around the maypole and sing songs. We eat herring, potatoes and strawberries with cream. The children pick flowers and make wreaths. It is the longest day of the year and the sun almost never sets. We celebrate with family and friends outside all day.",
	},
	{
		Title: "Vintern",
		Sv:    "Det snöar ute och det är kallt. Jag tar på mig en varm jacka och vantar. Barnen bygger en snögubbe i trädgården och kastar snöbollar. Vi dricker varm choklad inne vid brasan. Vintern i Sverige är lång men vacker. Allt är vitt och tyst. Jag gillar att gå långa promenader i snön.",
		En:    "It is snowing outside and it is cold. I put on a warm jacket and mittens. The children build a snowman in the garden and throw snowballs. We drink hot chocolate inside by the fireplace. Winter in Sweden is long but beautiful. Everything is white and quiet. I like to take long walks in the snow.",
	},
	{
		Title: "Att resa med tåg",
		Sv:    "Jag ska resa till Göteborg med tåget. Jag köper en biljett på stationen och hittar min plats. Tåget avgår klockan nio och resan tar tre timmar. Jag tittar ut genom fönstret på de svenska landskapen. Det är vackert med gröna fält och röda stugor. Jag läser en bok och dricker kaffe från restaurangvagnen.",
		En:    "I am going to travel to Gothenburg by train. I buy a ticket at the station and find my seat. The train departs at nine and the journey takes three hours. I look out the window at the Swedish landscapes. It is beautiful with green fields and red cottages. I read a book and drink coffee from the dining car.",
	},
	{
		Title: "Fika",
		Sv:    "I Sverige är fika mycket viktigt. Man tar en paus och dricker kaffe med vänner eller kollegor. Det finns alltid något gott att äta, som en bulle eller en kaka. Fika är mer än bara kaffe, det är en social tradition. Man pratar om livet och tar det lugnt. Många svenskar fikar minst två gånger om dagen.",
		En:    "In Sweden, fika is very important. You take a break and drink coffee with friends or colleagues. There is always something tasty to eat, like a bun or a cookie. Fika is more than just coffee, it is a social tradition. You talk about life and take it easy. Many Swedes have fika at least twice a day.",
	},
	{
		Title: "Mitt hem",
		Sv:    "Jag bor i en liten lägenhet i Malmö. Den har två rum, ett kök och ett badrum. Mitt favoritrum är vardagsrummet. Där har jag en soffa och en bokhylla full med böcker. Köket är litet men jag gillar att laga mat där. Från mitt fönster kan jag se Öresundsbron. Jag trivs bra i min lägenhet.",
		En:    "I live in a small apartment in Malmö. It has two rooms, a kitchen and a bathroom. My favorite room is the living room. There I have a sofa and a bookshelf full of books. The kitchen is small but I like to cook there. From my window I can see the Öresund Bridge. I am happy in my apartment.",
	},
	{
		Title: "På stranden",
		Sv:    "Sommaren i Sverige är kort men underbar. Vi åker till stranden och badar i havet. Vattnet är kallt men det känns skönt när solen skiner. Barnen bygger sandslott och vi äter glass. Vi lägger oss på filtar och njuter av solen. Måsarna flyger över oss och vi hör vågornas ljud. Det är en perfekt sommardag.",
		En:    "Summer in Sweden is short but wonderful. We go to the beach and swim in the sea. The water is cold but it feels nice when the sun is shining. The children build sandcastles and we eat ice cream. We lie on blankets and enjoy the sun. The seagulls fly above us and we hear the sound of the waves. It is a perfect summer day.",
	},
	{
		Title: "Skolan",
		Sv:    "I Sverige börjar barnen skolan när de är sex år. De lär sig läsa, skriva och räkna. Läraren är snäll och hjälper alla elever. På rasten leker barnen ute på skolgården. De äter lunch i matsalen varje dag. Maten är gratis för alla barn i Sverige. Efter skolan går många barn till fritids.",
		En:    "In Sweden, children start school when they are six years old. They learn to read, write and count. The teacher is kind and helps all the students. During recess the children play outside in the schoolyard. They eat lunch in the dining hall every day. The food is free for all children in Sweden. After school many children go to after-school care.",
	},
	{
		Title: "Lucia",
		Sv:    "Den trettonde december firar vi Lucia i Sverige. Barnen sjunger vackra sånger och bär vita kläder. Lucia har ljus i håret och leder processionen. Vi äter lussekatter och pepparkakor och dricker glögg. Det är en fin tradition i mörkaste december. Lucia bringer ljus i vintermörkret och alla känner sig glada.",
		En:    "On the thirteenth of December we celebrate Lucia in Sweden. The children sing beautiful songs and wear white clothes. Lucia has candles in her hair and leads the procession. We eat saffron buns and gingerbread cookies and drink mulled wine. It is a lovely tradition in the darkest December. Lucia brings light in the winter darkness and everyone feels happy.",
	},
	{
		Title: "Husdjur",
		Sv:    "Vi har en hund som heter Molly. Hon är brun och vit och mycket lekfull. Varje morgon går vi ut och promenerar med henne i skogen. Hon älskar att springa och leka med andra hundar. På kvällen ligger hon i soffan bredvid oss. Hon är en del av familjen och vi älskar henne mycket.",
		En:    "We have a dog called Molly. She is brown and white and very playful. Every morning we go out and walk with her in the forest. She loves to run and play with other dogs. In the evening she lies on the sofa next to us. She is part of the family and we love her very much.",
	},
	{
		Title: "Att laga mat",
		Sv:    "I kväll ska jag laga köttbullar med potatis och lingonsylt. Först steker jag köttbullarna i en stekpanna tills de är gyllene. Sedan kokar jag potatisen och gör en gräddsås. Det är en klassisk svensk middag som hela familjen gillar. Vi dukar bordet och tänder ljus. Det luktar underbart i hela köket.",
		En:    "Tonight I am going to make meatballs with potatoes and lingonberry jam. First I fry the meatballs in a frying pan until they are golden. Then I boil the potatoes and make a cream sauce. It is a classic Swedish dinner that the whole family likes. We set the table and light candles. It smells wonderful in the whole kitchen.",
	},
	{
		Title: "Väder i Sverige",
		Sv:    "Vädret i Sverige ändras mycket genom året. På våren blommar träden och det blir varmare. Sommaren är ljus med långa dagar och vita nätter. Hösten är färgglad med röda och gula löv som faller från träden. Vintern är mörk och kall med snö. Svenskar pratar ofta om vädret och det påverkar humöret.",
		En:    "The weather in Sweden changes a lot throughout the year. In spring the trees bloom and it gets warmer. Summer is bright with long days and white nights. Autumn is colorful with red and yellow leaves falling from the trees. Winter is dark and cold with snow. Swedes often talk about the weather and it affects the mood.",
	},
	{
		Title: "Fotboll",
		Sv:    "Min bror spelar fotboll varje lördag. Hans lag heter Blåvitt och de spelar på en stor plan. Jag brukar gå och titta på matcherna med min pappa. Det är roligt när de gör mål och alla jublar. Efter matchen dricker spelarna vatten och pratar om hur det gick. Min bror drömmer om att bli proffs.",
		En:    "My brother plays football every Saturday. His team is called Blåvitt and they play on a big field. I usually go and watch the matches with my dad. It is fun when they score a goal and everyone cheers. After the match the players drink water and talk about how it went. My brother dreams of becoming a professional.",
	},
	{
		Title: "Biblioteket",
		Sv:    "Jag går ofta till biblioteket. Där kan man låna böcker gratis. Jag gillar att läsa romaner och deckare. Det är tyst och lugnt på biblioteket. Ibland sitter jag där och studerar i flera timmar. Bibliotekarien hjälper mig att hitta nya böcker. Det finns också tidningar och datorer man kan använda.",
		En:    "I often go to the library. There you can borrow books for free. I like to read novels and crime fiction. It is quiet and calm at the library. Sometimes I sit there and study for several hours. The librarian helps me find new books. There are also newspapers and computers you can use.",
	},
	{
		Title: "Kräftskiva",
		Sv:    "I augusti har vi kräftskiva. Vi sitter ute i trädgården och äter kräftor med dill och bröd. Vi sjunger snapsvisor och bär roliga hattar och haklappar. Det är en gammal svensk tradition som alla älskar. Alla skrattar och har roligt tillsammans. Vi sitter ute sent på kvällen under lyktorna i det varma mörkret.",
		En:    "In August we have a crayfish party. We sit outside in the garden and eat crayfish with dill and bread. We sing drinking songs and wear funny hats and bibs. It is an old Swedish tradition that everyone loves. Everyone laughs and has fun together. We sit outside late in the evening under the lanterns in the warm darkness.",
	},
	{
		Title: "Pendla till jobbet",
		Sv:    "Jag pendlar till jobbet varje dag. Först går jag till tunnelbanan och åker tre stationer. Sedan byter jag till bussen som tar mig till kontoret. Det tar ungefär fyrtio minuter. Jag lyssnar på musik eller en podd under resan. Ibland läser jag nyheter på telefonen. Det är skönt att ha en rutin.",
		En:    "I commute to work every day. First I walk to the subway and ride three stations. Then I transfer to the bus that takes me to the office. It takes about forty minutes. I listen to music or a podcast during the trip. Sometimes I read the news on my phone. It is nice to have a routine.",
	},
	{
		Title: "Allemansrätten",
		Sv:    "I Sverige har vi allemansrätten. Det betyder att alla får vandra fritt i naturen. Man kan plocka bär och svamp i skogen utan att fråga. Man får tälta en natt nästan var som helst. Det är viktigt att inte skräpa ner och att vara rädd om naturen. Allemansrätten gör att alla kan njuta av den svenska naturen.",
		En:    "In Sweden we have the right of public access. It means that everyone may walk freely in nature. You can pick berries and mushrooms in the forest without asking. You may camp for one night almost anywhere. It is important not to litter and to take care of nature. The right of public access means that everyone can enjoy Swedish nature.",
	},
	{
		Title: "En dag på jobbet",
		Sv:    "Jag jobbar som lärare på en skola. Jag undervisar i svenska och engelska. Mina elever är snälla och nyfikna. Vi har rast klockan tio och äter lunch klockan tolv. På eftermiddagen har vi ofta grupparbeten. Jag rättar prov på kvällen hemma. Jag trivs bra med mitt jobb och älskar att se eleverna lära sig.",
		En:    "I work as a teacher at a school. I teach Swedish and English. My students are kind and curious. We have a break at ten and eat lunch at twelve. In the afternoon we often have group work. I grade tests in the evening at home. I enjoy my job and love to see the students learn.",
	},
	{
		Title: "Att handla kläder",
		Sv:    "Jag behöver nya kläder till hösten. Jag går till en affär i stan och provar en jacka. Den blå jackan passar bra och är inte för dyr. Jag köper också en varm tröja och ett par handskar. Expediten är hjälpsam och föreslår en fin halsduk som passar till jackan. Jag är nöjd med mina nya kläder.",
		En:    "I need new clothes for autumn. I go to a shop in town and try on a jacket. The blue jacket fits well and is not too expensive. I also buy a warm sweater and a pair of gloves. The shop assistant is helpful and suggests a nice scarf that goes with the jacket. I am pleased with my new clothes.",
	},
	{
		Title: "Dalarna",
		Sv:    "Dalarna ligger i mitten av Sverige. Det är känt för röda stugor och dalahästar. Många turister besöker Dalarna på sommaren för att vandra och bada i sjöarna. Naturen är vacker med djupa skogar och höga berg. Folk dansar runt midsommarstången och spelar folkmusik. Dalarna kallas ibland Sveriges hjärta.",
		En:    "Dalarna is located in the middle of Sweden. It is known for red cottages and Dala horses. Many tourists visit Dalarna in the summer to hike and swim in the lakes. The nature is beautiful with deep forests and tall mountains. People dance around the maypole and play folk music. Dalarna is sometimes called the heart of Sweden.",
	},
	{
		Title: "Fredagsmys",
		Sv:    "På fredagar har vi fredagsmys hemma. Vi lagar tacos och dricker läsk. Hela familjen sitter i soffan och tittar på en film tillsammans. Vi äter godis och chips och myser under filtar. Det är det bästa med veckan. Alla ser fram emot fredagen. Barnen får välja filmen och vi har det jättebra.",
		En:    "On Fridays we have cozy Friday at home. We make tacos and drink soda. The whole family sits on the sofa and watches a movie together. We eat candy and chips and relax under blankets. It is the best part of the week. Everyone looks forward to Friday. The children get to choose the movie and we have a great time.",
	},
	{
		Title: "Söka jobb",
		Sv:    "Jag söker ett nytt jobb. Jag skriver mitt CV och ett personligt brev. Jag skickar min ansökan till flera företag. Nästa vecka har jag en intervju. Jag är nervös men också glad och förbereder mig noga. Jag läser om företaget och övar på vanliga intervjufrågor. Jag hoppas att det går bra.",
		En:    "I am looking for a new job. I write my CV and a cover letter. I send my application to several companies. Next week I have an interview. I am nervous but also happy and prepare carefully. I read about the company and practice common interview questions. I hope it goes well.",
	},
	{
		Title: "Stockholms tunnelbana",
		Sv:    "Stockholms tunnelbana är känd för sin konst. Många stationer har målningar och skulpturer på väggarna. Det kallas världens längsta konstutställning. Jag tycker om att titta på konsten när jag väntar på tåget. Varje station ser annorlunda ut med olika färger och former. Det gör att pendla till jobbet blir lite roligare.",
		En:    "Stockholm's subway is known for its art. Many stations have paintings and sculptures on the walls. It is called the world's longest art exhibition. I like to look at the art when I am waiting for the train. Every station looks different with various colors and shapes. It makes commuting to work a little more fun.",
	},
	{
		Title: "Semester i fjällen",
		Sv:    "I vinter ska vi åka till fjällen. Vi ska bo i en stuga och åka skidor varje dag. Barnen vill åka pulka och bygga snögrottor. På kvällen tänder vi en brasa och dricker varm choklad. Vi spelar brädspel och berättar historier. Luften är kall och frisk och utsikten är fantastisk. Det blir en underbar semester.",
		En:    "This winter we are going to the mountains. We will stay in a cabin and ski every day. The children want to go sledding and build snow caves. In the evening we light a fire and drink hot chocolate. We play board games and tell stories. The air is cold and fresh and the view is fantastic. It will be a wonderful vacation.",
	},
}
