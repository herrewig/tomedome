import { markupAbilityAnswer } from "./data.js";

test('markupAbilityAnswer formats correctly', () => {
    const question = {
        abilityName: 'Blink',
        answer: {
            'description': 'teleport up to 600 units away',
            'attributes': ['foo', 'bar']
        }
    };
    const result = markupAbilityAnswer(question);
    expect(result).toBe('<p class="desc"><strong>Blink: </strong>teleport up to 600 units away</p><p class="attr">foo<br>bar</p>');
});