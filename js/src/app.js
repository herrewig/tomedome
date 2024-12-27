import { fetchQuiz } from './data.js';
import { setUpQuiz, renderQuestions,  } from './render.js';

let questionIndex = 0;
let isPrompting = false;

document.addEventListener('DOMContentLoaded', async () => {
    const container = document.getElementById("quiz-questions-container");
    const quiz = await fetchQuiz();
    const nextButton = document.getElementById("next-question");

    setUpQuiz(quiz);

    document.getElementById("next-question").addEventListener("click", () => {
        if (!isPrompting) {
            nextButton.textContent = "Show answer";
        } else {
            nextButton.textContent = "Next question";
        }

        renderQuestions(container, quiz.questions, questionIndex, isPrompting);

        if (isPrompting) {
            questionIndex++;
        }
        isPrompting = !isPrompting;

        // Disable next button if we're at the end of the quiz
        if (questionIndex >= quiz.questions.length) {
            nextButton.toggleAttribute("disabled");
        }
    });
});