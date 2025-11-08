package vocab_testing

import lsproto "vocab/lsp"

func FilterDiag(diags []lsproto.Diagnostic, sev lsproto.DiagnosticsSeverity) []lsproto.Diagnostic {
	var collected []lsproto.Diagnostic
	for _, diag := range diags {
		if diag.Severity == sev {
			collected = append(collected, diag)
		}
	}
	return collected
}
