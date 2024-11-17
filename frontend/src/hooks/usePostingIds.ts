import { useState, useEffect, useMemo } from 'react';
import { useBusinessInfo } from './useBusinessInfo.ts';
import { usePostingInfo } from './usePostingInfo.ts';

export function usePostingIds(): [string, number][] {
    const postingInfo = usePostingInfo();
    return useMemo(() => postingInfo.map((posting) => [posting.business_id, posting.id]), [postingInfo]);
}