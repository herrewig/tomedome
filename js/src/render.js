import { markupAbilityAnswer } from './data.js';

export function setUpQuiz(quiz) {
    let heroTitleElement = document.getElementById("hero-title");
    heroTitleElement.textContent = quiz.displayName;

    let heroImgElement = document.getElementById("hero-img");
    heroImgElement.src = `https://cdn.cloudflare.steamstatic.com/apps/dota2/images/dota_react/heroes/${quiz.shortName}.png`;
}

export function renderQuestions(container, questions, questionIndex, isPrompting) {
    let index = 0;
    container.innerHTML = "";

    while (index <= questionIndex) {
        const shouldAnswer = isPrompting || index < questionIndex;
        let prompt = null;
        let answer = null;

        // Create question
        const quizQuestion = document.createElement("div");
        quizQuestion.classList.add("quiz-question");

        // Create prompt
        prompt = document.createElement("div");
        prompt.classList.add("quiz-prompt");
        const promptP = document.createElement("p");
        promptP.classList.add("prompt");
        promptP.innerHTML = `<strong>${questions[index].prompt}</strong>`;
        prompt.appendChild(promptP);
        quizQuestion.appendChild(prompt);

        // Answer
        if (shouldAnswer) {
            answer = document.createElement("div");
            answer.classList.add("quiz-answer");
            const answerP = document.createElement("p");
            answerP.classList.add("answer-p");
            answerP.innerHTML = markupAbilityAnswer(questions[index]);
            answer.appendChild(answerP);
            quizQuestion.appendChild(answer);
        }

        // Append to container
        container.appendChild(quizQuestion);

        // Fade in newest element
        if (index == questionIndex) {
            let hidden = shouldAnswer ? answer : prompt;
            hidden.classList.add("hidden");
            setTimeout(() => {
                hidden.classList.remove("hidden");
                hidden.classList.add("fade-in");
            }, 0);
        }

        index++;
    }
}