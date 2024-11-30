import { useState, useEffect, useMemo } from 'react';
import { useBusinessInfo } from './useBusinessInfo.ts';
import { usePostingInfo } from './usePostingInfo.ts';

export function usePostingIds(): [string, number][] {
    const { data } = usePostingInfo();
    return useMemo(() => data.map((posting) => [posting.business_id, posting.id]), [data]);
}

export function usePostingNames(): Map<number, string> {
    const { data } = usePostingInfo();
    const hashmap = new Map<number, string>();
    data.map((posting) => hashmap.set(posting.id, posting.title));
    return hashmap;
}
