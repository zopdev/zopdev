export class HttpErrors {
  constructor(error = {}, status) {
    this.isHttpError = true;
    this.name = 'Http Error';
    this.status = status;
    this.message = error.message || 'An unknown error occurred.';
    this.details = error.details || error.detail || {};
  }
}

export const isHttpError = (resp) => resp && resp.isHttpError === true;
