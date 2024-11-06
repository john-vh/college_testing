import { useState, useEffect } from 'react';
import { useBusinessInfo } from './useBusinessInfo';

export function useBusinessIds(): string[] {
    return useBusinessInfo()?.map((business) => business.id) ?? [];
}