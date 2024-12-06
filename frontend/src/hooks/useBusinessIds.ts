import { useState, useEffect, useMemo } from 'react';
import { useBusinessInfo } from './useBusinessInfo.ts';

export function useBusinessIds(): string[] {
    const { businessInfo } = useBusinessInfo({});
    return useMemo(() => businessInfo.map((business) => business.id), [businessInfo]);
}