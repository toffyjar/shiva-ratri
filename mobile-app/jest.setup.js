// Silence the app's informational console.log calls (e.g. httpService's
// request logging) so test output only shows failures and warnings.
global.console.log = jest.fn();
