# Taste

## Workflow
- Prefers high effort as the default for sessions: comprehensive implementation with extensive testing and documentation. Confidence: 0.7
- Prefers asking clarifying questions before implementing whenever there are genuine design ambiguities, offering a single "Recommended" option per decision. Confidence: 0.8
- Prefers validating implementation with functional end-to-end testing (mock servers, throwaway local databases, `php -l`, following redirects) rather than relying only on unit or syntax checks. Confidence: 0.8
- Prefers the simplest tool that fits the requirement: favors a single self-contained file over a heavyweight multi-dependency alternative when it meets the need. Confidence: 0.7
- Prefers to retain control over environment teardown: objects to the agent auto-killing running test processes/servers or auto-deleting temp files; will handle fixes and cleanup themselves. Confidence: 0.5
