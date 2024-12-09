import { Button, Card, FormGroup, H2, H3, InputGroup, Intent, Spinner, Tag, TextArea } from "@blueprintjs/core";
import { Toaster, Position } from "@blueprintjs/core";
import React, { useMemo, useState } from "react";
import { AddPosting } from "./AddPosting.tsx";
import { usePostingInfo } from "../hooks/usePostingInfo.ts";
import { useActivatePosting } from "../hooks/useActivatePosting.ts";
import { useDeactivatePosting } from "../hooks/useDeactivatePosting.ts";
import { PostingInfo } from "../hooks/useAllPostings.ts";

const toaster = Toaster.create({
    position: Position.TOP,
});

interface PostingInfoProps {
    isAdmin: boolean;
}

export const PostingInfoPage = ({ isAdmin = false }: PostingInfoProps) => {
    const { data = [], business_map, error, fetchPostingInfo } = usePostingInfo({ isAdmin });
    const { activatePosting, loading: activateLoading } = useActivatePosting();
    const { deactivatePosting, loading: deactivateLoading } = useDeactivatePosting();
    const [showAddPosting, setShowAddPosting] = useState(false);

    const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

    const businessOptions = useMemo(
        () => [...new Set(data.map((post) => post.business_id))],
        [data]
    );

    const sortedData = useMemo(
        () => [...data].sort((a, b) => a.title.localeCompare(b.title)),
        [data]
    );

    const handleStatusChange = async (posting: PostingInfo) => {
        try {
            if (posting.status === "disabled") {
                await activatePosting(posting);
                toaster.show({
                    message: "Posting activated successfully",
                    intent: Intent.SUCCESS,
                });
            } else {
                await deactivatePosting(posting);
                toaster.show({
                    message: "Posting deactivated successfully",
                    intent: Intent.SUCCESS,
                });
            }
            await delay(100);
            await fetchPostingInfo();
        } catch (err) {
            toaster.show({
                message: `Error changing posting status: ${err instanceof Error ? err.message : 'Unknown error'}`,
                intent: Intent.DANGER,
            });
        }
    };

    const handleAddNewPosting = async () => {
        setShowAddPosting(true);
        try {
            console.log("oops too early");
            await fetchPostingInfo();
        } catch (err) {
            toaster.show({
                message: "Error refreshing posting list",
                intent: Intent.DANGER,
            });
        }
    };

    if (error) {
        return (
            <Card style={{ margin: "200px" }} elevation={2}>
                <H3 style={{ color: Intent.DANGER }}>Error loading postings</H3>
                <p>{error}</p>
            </Card>
        );
    }

    if (showAddPosting) {
        return (
            <AddPosting
                businesses={businessOptions}
                setPostingAdd={setShowAddPosting}
                fetchData={fetchPostingInfo}
            />
        );
    }

    return (
        <div className="content">
            <div className="Flex" style={{ justifyContent: "space-between", alignItems: "center" }}>
                <H2 style={{ marginBottom: "0px" }}>Manage Postings</H2>
                <Button
                    intent="primary"
                    large

                    onClick={() => handleAddNewPosting()}
                >
                    Add Posting
                </Button>
            </div>

            <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
                {sortedData.map((posting) => (
                    <Card
                        key={posting.id}
                        interactive={true}
                        elevation={2}
                        style={{ position: "relative" }}
                    >
                        {/* {(activateLoading || deactivateLoading) && (
                            <div
                                style={{
                                    position: "absolute",
                                    display: "flex",
                                    alignItems: "center",
                                    justifyContent: "center",
                                }}
                            >
                                <Spinner size={50} />
                            </div>
                        )} */}

                        <div className="Flex" style={{ justifyContent: "space-between" }}>
                            <H3>Posting Information</H3>
                            <Tag round intent={posting.status === "active" ? Intent.SUCCESS : Intent.DANGER}>{posting.status.toLocaleUpperCase()}</Tag>
                        </div>

                        <FormGroup inline label="Title" labelFor={`title-${posting.id}`}>
                            <InputGroup
                                id={`title-${posting.id}`}
                                value={posting.title}
                                readOnly
                                fill
                            />
                        </FormGroup>

                        <FormGroup inline label="Business" labelFor={`business-${posting.id}`}>
                            <InputGroup
                                id={`business-${posting.id}`}
                                value={business_map.get(posting.business_id)?.name}
                                readOnly
                                fill
                            />
                        </FormGroup>

                        <FormGroup inline label="ID" labelFor={`id-${posting.id}`}>
                            <InputGroup
                                id={`id-${posting.id}`}
                                value={posting.id.toString()}
                                readOnly
                                fill
                            />
                        </FormGroup>

                        <FormGroup label="Description" labelFor={`desc-${posting.id}`}>
                            <TextArea
                                autoResize
                                asyncControl
                                id={`desc-${posting.id}`}
                                value={posting.desc}
                                readOnly
                                fill
                            />
                        </FormGroup>

                        <Button
                            intent={posting.status === "disabled" ? Intent.SUCCESS : Intent.DANGER}
                            onClick={() => handleStatusChange(posting)}
                            disabled={activateLoading || deactivateLoading}
                        >
                            {posting.status === "disabled" ? "Activate posting" : "Deactivate posting"}
                        </Button>
                    </Card>
                ))}
            </div>
        </div>
    );
};