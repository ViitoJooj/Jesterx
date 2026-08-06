export interface User {
  uuid: string;
  websiteUuid: string;
  imageUrl: string | null;
  name: string;
  email: string;
  role: string;
  password: string;
  cpf: string | null;
  githubOauth: boolean;
  googleOauth: boolean;
  appleOauth: boolean;
  updatedAt: string | null;
  createdAt: string;
}
