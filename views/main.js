const canvas = new fabric.Canvas("canvas");
let is_first = false
canvas.backgroundColor = "white";
const sleep = ms => new Promise(res => setTimeout(res, ms))

window.onload = function () {
    document.getElementById("draw").click();
    fetch("/api/oekaki")
        .then(res => res.json())
        .then(data => {
            const img_src = data.image
            if (img_src == "") {
                is_first = true
                console.log(is_first)
                document.getElementById("sorry").style.display = "block";
            } else {
                let odai = document.getElementById("odai")
                odai.src = img_src;
                odai.style.display = "block";
                document.getElementById("answer").style.display = "block";
                document.getElementById("answer_desc").style.display = "block";
            }
        })
}



function validate_hiragana() {
    const hiragana_regex = `[\u3041-\u3096]*`
    let answer = document.getElementById("answer").value
    let next_answer = document.getElementById("next_answer").value

    let submit_button = document.getElementById("submit_button")

    if ((answer.match(hiragana_regex) || is_first) && next_answer.match(hiragana_regex)) {
        submit_button.disabled = false;
        console.log("false")
    } else {
        submit_button.disabled = true;
        console.log("true")
    }
}

document.getElementById("answer").addEventListener('change', function (e) { validate_hiragana() });
document.getElementById("next_answer").addEventListener('change', function (e) { validate_hiragana() });

async function do_alert(message) {
    document.getElementById("alert").innerHTML = message;
    document.getElementById("alert").style.display = "block";
    await sleep(3000)
    document.getElementById("alert").style.display = "none";
}

async function post_data() {
    document.getElementById("submit_button").disabled = true

    let data = {
        image: document.getElementById("canvas").toDataURL(),
        answer: document.getElementById("answer").value,
        next_answer: document.getElementById("next_answer").value,
    }

    const res_raw = await fetch('/api/oekaki', {
        method: 'post',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    const res = await res_raw.json()
    if (res.success) {
        if (res.correct && !is_first) {
            await do_alert("あなたの回答「" + document.getElementById("answer").value + "」は正解でした！<span class=\"inline-block\">あなたのお題を送信しました。</span>")
        } else if (!res.correct && !is_first) {
            await do_alert("あなたの回答「" + document.getElementById("answer").value + "」は不正解でした...<span class=\"inline-block\">あなたのお題を送信しました。</span>")
        } else {
            await do_alert("あなたのお題を送信しました。")
        }
        location.reload();
    } else {
        let additional = ""
        if (res.error == "Too small input image") {
            additional = "もう少し、何か描いてあげてください。"
        }
        if (res.error.includes("non hiragana character")) {
            additional = "回答は、すべてひらがなで行なってください。"
        }
        if (res.error.includes("answer is not given")) {
            additional = "なにか回答を入力してください。"
        }
        document.getElementById("submit_button").disabled = false;
        await do_alert("送信に失敗しました。" + additional)
    }
}

document.getElementById("draw").addEventListener("click", function () {
    canvas.freeDrawingBrush = new fabric.PencilBrush(canvas);
    canvas.freeDrawingBrush.width = 5;
    canvas.freeDrawingBrush.color = "black";
    canvas.isDrawingMode = true;
});

document.getElementById("erase").addEventListener("click", function () {
    canvas.freeDrawingBrush = new fabric.PencilBrush(canvas);
    canvas.freeDrawingBrush.width = 20;
    canvas.freeDrawingBrush.color = "white";
    canvas.isDrawingMode = true;
});
