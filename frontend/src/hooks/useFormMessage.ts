import { useState } from 'react';

export type MessageStatus = 'success' | 'error' | null;

export interface UseFormMessageReturn {
  message: string;
  status: MessageStatus;
  setMessage: (message: string) => void;
  setStatus: (status: MessageStatus) => void;
  setError: (error: string) => void;
  setSuccess: (success: string) => void;
  clearMessage: () => void;
}

export function useFormMessage(): UseFormMessageReturn {
  const [message, setMessage] = useState<string>('');
  const [status, setStatus] = useState<MessageStatus>(null);

  const setError = (error: string) => {
    setMessage(error);
    setStatus('error');
  };

  const setSuccess = (success: string) => {
    setMessage(success);
    setStatus('success');
  };

  const clearMessage = () => {
    setMessage('');
    setStatus(null);
  };

  return {
    message,
    status,
    setMessage,
    setStatus,
    setError,
    setSuccess,
    clearMessage,
  };
}
