import { markupAbilityAnswer } from "./data.js";

test('markupAbilityAnswer formats correctly', () => {
    const question = { abilityName: 'Blink', answer: 'teleport up to 600 units away' };
    const result = markupAbilityAnswer(question);
    expect(result).toBe('<strong>Blink: </strong>teleport up to 600 units away');
});