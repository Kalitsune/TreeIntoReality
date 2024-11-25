use iocraft::prelude::*;

#[derive(Default, Props)]
struct TreeInputProps<'a> {
    value_out: Option<&'a mut String>,
}

#[component]
fn TreeInput<'a>(props: &mut TreeInputProps<'a>, mut hooks: Hooks) -> impl Into<AnyElement<'static>> {
    let mut system = hooks.use_context_mut::<SystemContext>();

    let mut value = hooks.use_state(|| "".to_string());

    let mut should_submit = hooks.use_state(|| false);
    let mut should_stop = hooks.use_state(|| false);

    hooks.use_terminal_events(move |event| match event {
        TerminalEvent::Key(KeyEvent { code, modifiers, kind, .. }) if kind != KeyEventKind::Release => {
            match code {
                KeyCode::Enter => should_submit.set(true && modifiers.contains(KeyModifiers::SHIFT)),
                KeyCode::Esc => should_stop.set(true),
                _ => {}
            }
        }
        _ => {}
    });

    if should_submit.get() {
        if let Some(value_out) = props.value_out.as_mut() {
            **value_out = value.to_string();
        }
        system.exit();
    }
    element! {
        Box(
            flex_direction: FlexDirection::Column,
            border_style: BorderStyle::Round,
            border_color: Color::Blue,
            padding_top: 1,
            padding_bottom: 1,
            padding_right: 1,
            padding_left: 1,
        ) {
            Text(
                color: Color::Blue,
                content: "Paste your tree output here:"
            )
            Box(
                background_color: Color::Black,
                min_height: 5,
            ) {
                TextInput(
                    has_focus: true,
                    value: value.to_string(),
                    on_change: move |new_value| value.set(new_value),
                )
            }
            Text(
                color: Color::DarkGrey
                ,content: "Validate: [Ctrl] + [S]"
            )
        }
    }
}
pub fn interactive() {
    let mut tree = String::new();

    // Welcome screen
    smol::block_on(
        element! {
            Box(
                flex_direction: FlexDirection::Column,
                width: 100,
            ){
                Box(
                    border_style: BorderStyle::Round,
                    border_color: Color::Green,
                    flex_direction: FlexDirection::Column,
                    justify_content: JustifyContent::Center,
                ) {
                    Text(content: " ║ Tree Into Reality  ")
                    Text(content: " ╚══════════════════════════════════════════════════════════════════════════════════════════════⇒")
                }
                TreeInput(value_out: &mut tree)
            }
        }
        .render_loop(),
    ).unwrap();

    if !tree.is_empty() {
        println!("{}", tree.to_string())
    }

}