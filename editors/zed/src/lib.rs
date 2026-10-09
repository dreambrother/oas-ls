use zed_extension_api::{self as zed, settings::LspSettings, LanguageServerId, Result};

const LANGUAGE_SERVER_ID: &str = "oas-ls";
const BINARY_NAME: &str = "oas-ls";

struct OasLsExtension;

impl OasLsExtension {
    fn binary_path(
        &mut self,
        binary_settings: Option<&zed::settings::CommandSettings>,
        worktree: &zed::Worktree,
    ) -> Result<String> {
        if let Some(path) = binary_settings.and_then(|binary| binary.path.clone()) {
            return Ok(path);
        }

        if let Some(path) = worktree.which(BINARY_NAME) {
            return Ok(path);
        }

        if std::fs::metadata(BINARY_NAME).map(|stat| stat.is_file()).unwrap_or(false) {
            return Ok(BINARY_NAME.to_string());
        }

        Err(format!(
            "{BINARY_NAME} was not found. Build it with `go install ./cmd/{BINARY_NAME}` from the oas-ls repository (or `go install github.com/dreambrother/oas-ls/cmd/{BINARY_NAME}@latest`) so it lands on PATH (for example ~/.go/bin), or set `lsp.{LANGUAGE_SERVER_ID}.binary.path` to the built binary in your Zed settings."
        ))
    }
}

impl zed::Extension for OasLsExtension {
    fn new() -> Self {
        Self
    }

    fn language_server_command(
        &mut self,
        _language_server_id: &LanguageServerId,
        worktree: &zed::Worktree,
    ) -> Result<zed::Command> {
        let binary_settings = LspSettings::for_worktree(LANGUAGE_SERVER_ID, worktree)
            .ok()
            .and_then(|settings| settings.binary);
        let command = self.binary_path(binary_settings.as_ref(), worktree)?;
        let args = binary_settings
            .as_ref()
            .and_then(|binary| binary.arguments.clone())
            .unwrap_or_default();
        let env = binary_settings
            .as_ref()
            .and_then(|binary| binary.env.clone())
            .map(|env| env.into_iter().collect::<Vec<_>>())
            .unwrap_or_default();
        Ok(zed::Command { command, args, env })
    }

    fn language_server_initialization_options(
        &mut self,
        server_id: &LanguageServerId,
        worktree: &zed::Worktree,
    ) -> Result<Option<zed_extension_api::serde_json::Value>> {
        Ok(LspSettings::for_worktree(server_id.as_ref(), worktree)
            .ok()
            .and_then(|settings| settings.initialization_options))
    }

    fn language_server_workspace_configuration(
        &mut self,
        server_id: &LanguageServerId,
        worktree: &zed::Worktree,
    ) -> Result<Option<zed_extension_api::serde_json::Value>> {
        Ok(LspSettings::for_worktree(server_id.as_ref(), worktree)
            .ok()
            .and_then(|settings| settings.settings))
    }
}

zed::register_extension!(OasLsExtension);
