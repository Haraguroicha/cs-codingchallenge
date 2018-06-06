const topic = document.querySelector("#topic");
const pager = document.querySelector("#pager");
const topicBox = document.querySelector("#topics");
const calcLeft = () => {
    document.querySelector("#left").innerHTML = `Left ${config.MaximumTopicLength - topic.value.length} characters.`;
};
calcLeft();
const createPager = (pages) => {
    pager.innerHTML = "";
    for (var p = 1; p <= pages.lastPage; p++) {
    var pg = document.createElement("option");
    pg.value = p;
    pg.innerHTML = p;
    if (p == pages.currentPage) {
        pg.selected = true;
    }
    pager.appendChild(pg);
    }
};
const createTopics = (topics) => {
    Array.from(document.querySelectorAll(".vote"))
        .map((v) => v.removeEventListener("click", doVote, false))
    topicBox.innerHTML = topics.map((t) => `<div class="col col-md-6">
        <div class="box">
            <span class="topic-id">#${t.topicId}</span>
            <div class="votes">
                <div class="vote up-vote">
                    <i class="fa fa-arrow-up fa-2x" aria-hidden="true"></i>
                    ${t.votes.upVotes}
                </div>
                <div class="vote down-vote">
                    <i class="fa fa-arrow-down fa-2x" aria-hidden="true"></i>
                    ${t.votes.downVotes}
                </div>
                Sum: ${t.votes.sumVotes}
            </div>
            <hr />
            <span class="topic-title">${t.topicTitle}</span>
        </div>
    </div>`
    ).join('\n');
    Array.from(document.querySelectorAll(".vote"))
        .map((v) => v.addEventListener("click", doVote, false))
}
const makeTopics = (topics) => {
    if (topics.success) {
        createPager(topics.pages);
        createTopics(topics.data);
        topic.value = "";
        calcLeft();
    } else {
        throw topics.message;
    }
};
const getTopics = (page) => {
    page = page || '1';
    fetch(`/api/getTopics/${page}`)
        .then((response) => response.json())
        .then(makeTopics)
        .catch((e) => {
            console.error(e);
            alert(e);
        });
};
const newTopic = (title) => {
    fetch("/api/newTopic", {
        body: JSON.stringify({"topicTitle": title}),
        method: 'POST'
    })
        .then((response) => response.json())
        .then(makeTopics)
        .catch((e) => {
            console.error(e);
            alert(e);
        })
        .then(() => {
            getTopics(Math.max(pager.selectedIndex, 0));
        });
};
const doVote = (e) => {
    const target = e.target;
    const icon = target.querySelector("i") || target;
    const votes = icon.parentElement;
    const box = votes.parentElement.parentElement;
    const topicId = box.querySelector(".topic-id").innerText.slice(1);
    const vote = votes.className.match(/(\w+-vote)/gi).pop().replace(/-(.)/gi, (a,b) => b.toUpperCase());
    const page = Math.max(pager.selectedIndex, 0) + 1;
    fetch(`/api/${vote}/${topicId}/${page}`, {
        method: 'POST'
    })
        .then((response) => response.json())
        .then(makeTopics)
        .catch((e) => {
            console.error(e);
            alert(e);
        });
};
topic.addEventListener("keypress", (e) => setTimeout(calcLeft, 100), false);
topic.addEventListener("keyup", (e) => setTimeout(calcLeft, 100), false);
topic.addEventListener("keydown", (e) => setTimeout(calcLeft, 100), false);
pager.addEventListener("change", (e) => getTopics(Math.max(pager.selectedIndex, 0) + 1), false);
document.querySelector(".topic-submit")
    .addEventListener("click", () => newTopic(topic.value), false);

getTopics();
