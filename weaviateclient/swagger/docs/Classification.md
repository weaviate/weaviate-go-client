# Classification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | ID to uniquely identify this classification run | [optional] [default to null]
**Class** | **string** | class (name) which is used in this classification | [optional] [default to null]
**ClassifyProperties** | **[]string** | which ref-property to set as part of the classification | [optional] [default to null]
**BasedOnProperties** | **[]string** | base the text-based classification on these fields (of type text) | [optional] [default to null]
**Status** | **string** | status of this classification | [optional] [default to null]
**Meta** | [***ClassificationMeta**](ClassificationMeta.md) | additional meta information about the classification | [optional] [default to null]
**Type_** | **string** | which algorythim to use for classifications | [optional] [default to null]
**K** | **int32** | k-value when using k-Neareast-Neighbor | [optional] [default to 3]
**InformationGainCutoffPercentile** | **int32** | Only available on type&#x3D;contextual. All words in a source corpus are ranked by their information gain against the possible target objects. A cutoff percentile of 40 implies that the top 40% are used and the bottom 60% are cut-off. | [optional] [default to 30]
**InformationGainMaximumBoost** | **int32** | Only available on type&#x3D;contextual. Words in a corpus will receive an additional boost based on how high they are ranked according to information gain. Setting this value to 3 implies that the top-ranked word will be ranked 3 times as high as the bottom ranked word. The curve in between is logarithmic. A maximum boost of 1 implies that no boosting occurs. | [optional] [default to 3]
**TfidfCutoffPercentile** | **int32** | Only available on type&#x3D;contextual. All words in a corpus are ranked by their tf-idf score. A cutoff percentile of 80 implies that the top 80% are used and the bottom 20% are cut-off. This is very effective to remove words that occur in almost all objects, such as filler and stop words. | [optional] [default to 80]
**MinimumUsableWords** | **int32** | Only available on type&#x3D;contextual. Both IG and tf-idf are mechanisms to remove words from the corpora. However, on very short corpora this could lead to a removal of all words, or all but a single word. This value guarantees that - regardless of tf-idf and IG score - always at least n words are used. | [optional] [default to 3]
**Error_** | **string** | error message if status &#x3D;&#x3D; failed | [optional] [default to null]
**SourceWhere** | [***WhereFilter**](WhereFilter.md) | limit the objects to be classified | [optional] [default to null]
**TrainingSetWhere** | [***WhereFilter**](WhereFilter.md) | Limit the training objects to be considered during the classification. Can only be used on types with explicit training sets, such as &#39;knn&#39; | [optional] [default to null]
**TargetWhere** | [***WhereFilter**](WhereFilter.md) | Limit the possible sources when using an algorithm which doesn&#39;t really on trainig data, e.g. &#39;contextual&#39;. When using an algorithm with a training set, such as &#39;knn&#39;, limit the training set instead | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


