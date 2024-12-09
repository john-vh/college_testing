import { useState, useEffect, useMemo } from 'react';
import { useBusinessInfo } from './useBusinessInfo.ts';
import { usePostingInfo } from './usePostingInfo.ts';

export function usePostingIds(isAdmin: boolean): [string, number][] {
    const { data } = usePostingInfo({ isAdmin });
    return useMemo(() => data.map((posting) => [posting.business_id, posting.id]), [data]);
}

export function usePostingNames(isAdmin: boolean): Map<number, string> {
    const { data } = usePostingInfo({ isAdmin });
    const hashmap = new Map<number, string>();
    data.map((posting) => hashmap.set(posting.id, posting.title));
    return hashmap;
}
