import base64

import cv2
import numpy as np
import os
import pickle
from flask import Flask, request, jsonify


def extract_features(image_path):
    image = cv2.imread(image_path)
    try:
        # Using SIFT to extract features
        gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
        sift = cv2.SIFT_create()
        keypoints, descriptors = sift.detectAndCompute(gray, None)
        descriptors = descriptors.flatten()
    except cv2.error as e:
        print('Error: ', e)
        return None

    return descriptors


def batch_extractor(images_path):
    files = [os.path.join(images_path, p) for p in sorted(os.listdir(images_path))]
    result = {}
    for f in files:
        print('Extracting features from image %s' % f)
        name = f.split('/')[-1].lower()
        result[name] = extract_features(f)
    return result


def build_inv_index(features):
    inv_index = {}
    for name in features.keys():
        for i in range(len(features[name])):
            if i not in inv_index:
                inv_index[i] = {}
            if name not in inv_index[i]:
                inv_index[i][name] = []
            inv_index[i][name].append(features[name][i])
    return inv_index


def search(query_path, index, top_k=10):
    query_features = extract_features(query_path)
    scores = {}
    for i in range(len(query_features)):
        if i not in index:
            continue
        for name in index[i]:
            if name not in scores:
                scores[name] = 0
            for feature in index[i][name]:
                scores[name] += np.dot(feature, query_features[i])
    results = sorted(scores.items(), key=lambda x: -x[1])[:top_k]
    return [r[0] for r in results]


def save_index(index, output_path):
    with open(output_path, 'wb') as f:
        pickle.dump(index, f)


def load_index(input_path):
    with open(input_path, 'rb') as f:
        return pickle.load(f)


print('loading...')
images_path = 'images'
index_path = 'index.pkl'
# features = batch_extractor(images_path)
# index = build_inv_index(features)
# save_index(index, index_path)

# query_path = 'test.jpg'
index = load_index(index_path)
print('done')

app = Flask(__name__)


@app.route('/search', methods=['POST'])
def search_images():
    # Check if image file is uploaded
    if 'image' not in request.files:
        return jsonify({'error': 'No image uploaded'}), 400

    # Extract features from uploaded image
    image = request.files['image']
    image_path = 'uploaded_image.jpg'
    image.save(image_path)
    query_results = Search(image_path)

    # Return top 10 search results as JSON
    # results = [{'image_path': path, 'similarity': score} for path, score in query_results]
    return jsonify({'results': query_results}), 200


def Search(query_path):
    results = search(query_path, index)
    print(results)
    return results


if __name__ == "__main__":
    app.run()
