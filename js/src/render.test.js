import { renderQuestions } from './render.js'; // Adjust the import to your file path

describe('renderQuestions', () => {
    let container;
    let questions;

    beforeEach(() => {
        // Create a mock container
        container = document.createElement('div');

        // Setup mock quiz data
        questions = [
            { prompt: 'What your Q?', abilityName: 'Stifling Dagger', answer: { 'description': 'blah blah', 'attributes': ['sup', 'yo'] } },
            { prompt: 'What your W?', abilityName: 'Phantom Strike', answer: { 'description': 'blah blah', 'attributes': ['nm', 'bra'] } },
            { prompt: 'What your E?', abilityName: 'Blur', answer: { 'description': 'blahhhhh', attributes: ['whats', 'good'] } },
        ];
    });

    it('should render first prompt of first question with no answer', () => {
        renderQuestions(container, questions, 0, false);

        // Check that two questions are rendered
        const quizQuestions = container.querySelectorAll('.quiz-question');
        expect(quizQuestions.length).toBe(1);

        // Check that the prompt is present in the first question
        const firstQuestionPrompt = quizQuestions[0].querySelector('.quiz-prompt p');
        expect(firstQuestionPrompt.innerHTML).toBe('<strong>What your Q?</strong>');

        // Check that the answer is not present in the first question
        const firstQuestionAnswer = quizQuestions[0].querySelector('.quiz-answer');
        expect(firstQuestionAnswer).toBeNull();
    });

    it('should render first question and answer', () => {
        renderQuestions(container, questions, 1, true);

        // Check that two questions are rendered
        const quizQuestions = container.querySelectorAll('.quiz-question');
        expect(quizQuestions.length).toBe(2);

        // Check that the prompt is present in the first question
        const firstQuestionPrompt = quizQuestions[1].querySelector('.quiz-prompt p');
        expect(firstQuestionPrompt.innerHTML).toBe('<strong>What your W?</strong>');

        // Check that the answer is present in the first question
        const firstQuestionAnswer = quizQuestions[1].querySelector('.quiz-answer p');
        expect(firstQuestionAnswer.innerHTML).toBe('<p class="desc"><strong>Phantom Strike: </strong>blah blah</p><p class="attr">nm<br>bra</p>');
    });

    it('should render all questions with prompt', () => {
        renderQuestions(container, questions, 2, false);

        // Check that two questions are rendered
        const quizQuestions = container.querySelectorAll('.quiz-question');
        expect(quizQuestions.length).toBe(3);

        // Check that the prompt is present in the first question
        const firstQuestionPrompt = quizQuestions[2].querySelector('.quiz-prompt p');
        expect(firstQuestionPrompt.innerHTML).toBe('<strong>What your E?</strong>');

        // Check that the answer is not present in the first question
        const firstQuestionAnswer = quizQuestions[2].querySelector('.quiz-answer');
        expect(firstQuestionAnswer).toBeNull();
    });

    it('should render all questions with answer', () => {
        renderQuestions(container, questions, 2, true);

        // Check that two questions are rendered
        const quizQuestions = container.querySelectorAll('.quiz-question');
        expect(quizQuestions.length).toBe(3);

        // Check that the prompt is present in the first question
        const firstQuestionPrompt = quizQuestions[2].querySelector('.quiz-prompt p');
        expect(firstQuestionPrompt.innerHTML).toBe('<strong>What your E?</strong>');

        // Check that the answer is not present in the first question
        const firstQuestionAnswer = quizQuestions[2].querySelector('.quiz-answer p');
        expect(firstQuestionAnswer.innerHTML).toBe('<p class="desc"><strong>Blur: </strong>blahhhhh</p><p class="attr">whats<br>good</p>');
    });
});